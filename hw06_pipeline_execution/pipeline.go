package hw06_pipeline_execution //nolint:golint,stylecheck
import (
	"fmt"
	"sync"
)

type Message struct {
	Order int32
	Value interface{}
}

type (
	In         = <-chan interface{}
	Out        = In
	Bi         = chan interface{}
	MessageBi  = chan Message
	MessageIn  = <-chan Message
	MessageOut = MessageIn
)

type Stage func(in In) (out Out)

// Так как последовательный pipeline было делать не интересно,
// добавляем возможность запуска нескольких stage в параллельных горутинах
// примерная схема работы
//                            x (stage0-0)				  (stage1-0)
// ->x(проставляем номера)	 				(соединяем)						(сортируем результаты по номеру)
//                            x (stage0-1)				   (stage1-1)
//stage0-0 - stage первой фазы в 0 параллельной горутине
// количество параллельно запущенных stage.
const parallelNumber = 20

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return createPipeline(stages, in, done, parallelNumber)
}

// Создаем pipeline с указанием количества параллельных stage.
func createPipeline(stages []Stage, in In, done In, numberInEachStages int) Out {
	nextChannel := convertMessage(in, done)
	for stageNumber, stage := range stages {
		nextChannel = runStageInParallel(stage, nextChannel, done, numberInEachStages, stageNumber)
	}
	return sortFromParallel(nextChannel, done)
}

// конвертируем сообщения
// необходимо для поддержания порядка сообщений при параллельной работе.
func convertMessage(in In, done In) MessageOut {
	out := make(MessageBi)
	go func() {
		defer close(out)
		var numberMessages int32 = 0
		for nextValue := range in {
			select {
			case <-done:
				fmt.Printf("convertMessage. close  - received done\n")
				return
			default:
				numberMessages++
				message := Message{numberMessages, nextValue}
				fmt.Printf("convertMessage. send value %#v\n", message)
				out <- message
			}
		}
	}()
	return out
}

// Запускает stage в параллеле.
func runStageInParallel(stage Stage, in MessageOut, done Out, numberInEachStages int, numberStage int) MessageOut {
	messageOut := make(MessageBi)
	go func() {
		// читается всеми читателями, поэтому количество указано, меньше на 1,
		// чтобы блокировалось, если не успеваем обработать - что-то типа backpressure
		funIn := make(MessageBi, numberInEachStages-1)

		// используется для завершения параллельных stage
		var wg sync.WaitGroup

		defer func(funIn MessageBi, messageOut MessageBi) {
			// порядок закрытия очень важен, сначала закрываем поток для stage
			close(funIn)
			// ждем прекращения работы stage
			wg.Wait()
			// закрываем выходной поток
			close(messageOut)
			fmt.Printf("runStageInParallel end, numberStage - %v\n", numberStage)
		}(funIn, messageOut)

		// ради чего все и затевалось - запуск stage в несколько горутин
		for numberGoroutine := 0; numberGoroutine < numberInEachStages; numberGoroutine++ {
			wg.Add(1)
			go runWrapper(&wg, stage, funIn, messageOut, numberStage, numberGoroutine)
		}

		for nextValue := range in {
			select {
			case <-done:
				fmt.Printf("runStageInParallel. close - received done, numberStage - %v\n", numberStage)
				return
			default:
				fmt.Printf("runStageInParallel new value: %#v, numberStage - %v\n", nextValue, numberStage)
				funIn <- nextValue
			}
		}
	}()
	return messageOut
}

// запуск обертки - требуется,  так как stage не умеет работать с сообщениями с номерами.
func runWrapper(wg *sync.WaitGroup, stage Stage, in MessageBi, out MessageBi, numberStage int, numberGoroutine int) {
	inStage := make(Bi)
	defer func(wg *sync.WaitGroup, inStage Bi) {
		wg.Done()
		close(inStage)
	}(wg, inStage)
	outStage := stage(inStage)
	for nextValue := range in {
		fmt.Printf("runWrapper newValue %#v, numberStage - %v, numberGoroutine %v\n",
			nextValue, numberStage, numberGoroutine)
		inStage <- nextValue.Value
		processedValue := <-outStage
		message := Message{nextValue.Order, processedValue}
		out <- message
		fmt.Printf("runWrapper sent out inStage %#v, numberStage - %v, numberGoroutine %v\n",
			message, numberStage, numberGoroutine)
	}
}

// сортируем получаемые сообщения, если сообщение еще не пришло, то ждем следующего.
func sortFromParallel(in MessageOut, done Out) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		results := make(map[int32]interface{})
		var nextNumber int32 = 1
		for nextValue := range in {
			select {
			case <-done:
				fmt.Printf("sortFromParallel. close  - received done\n")
				return
			default:
				trySendByOrder(nextValue, results, &nextNumber, out)
			}
		}
	}()
	return out
}

// получаем новый элемент, сохраняем в мапу, проверяем, есть ли по номеру что отправить
// если есть - отправляем, если нет, ждем следующего.
func trySendByOrder(nextValue Message, results map[int32]interface{}, nextNumber *int32, out Bi) {
	fmt.Printf("sortFromParallel new value: %#v\n", nextValue)
	results[nextValue.Order] = nextValue.Value
	for {
		foundValue, ok := results[*nextNumber]
		if !ok {
			break
		}
		fmt.Printf("sortFromParallel. sortFromParallel send: %#v\n", foundValue)
		out <- foundValue
		delete(results, *nextNumber)
		*nextNumber++
	}
}
