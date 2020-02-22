using System;
using System.IO;
using System.Text;
using System.Threading.Tasks;

namespace ParallelProgramming_lab1
{
    enum CharType
    {
        Letter,Digit,Other
    }

    class Program
    {
        static void ProducerAction(ref RingBuffer ringBuffer)
        {
            for (; ; )
            {
                var symbol = Console.ReadKey().KeyChar;

                ringBuffer.Empty.WaitOne();
                ringBuffer.Busy.WaitOne();

                ringBuffer.Buffer[ringBuffer.Head] = symbol;
                ringBuffer.Head = (ringBuffer.Head + 1) % ringBuffer.Buffer.Length;

                ringBuffer.Busy.Release();
                ringBuffer.Full.Release();
            }
        }

        static void ConsumerAction(ref RingBuffer ringBuffer, CharType charType)
        {
            for (; ; )
            {
                ringBuffer.Full.WaitOne();
                ringBuffer.Busy.WaitOne();

                char symbol = ringBuffer.Buffer[ringBuffer.Tail];

                switch (charType)
                {
                    case CharType.Letter:
                        if (char.IsLetter(symbol))
                        {
                            ringBuffer.Tail = (ringBuffer.Tail + 1) % ringBuffer.Buffer.Length;
                            
                            Console.ForegroundColor = ConsoleColor.Red;
                            Console.CursorLeft -= 1; ;
                            Console.Write(symbol);
                            Console.ResetColor();

                            ringBuffer.Busy.Release();
                            ringBuffer.Empty.Release();
                        }
                        else
                        {
                            ringBuffer.Full.Release();
                            ringBuffer.Busy.Release();
                        }
                        break;
                    case CharType.Digit:
                        if (char.IsDigit(symbol))
                        {
                            ringBuffer.Tail = (ringBuffer.Tail + 1) % ringBuffer.Buffer.Length;
                            
                            Console.ForegroundColor = ConsoleColor.Green;
                            Console.CursorLeft -= 1;
                            Console.Write(symbol);
                            Console.ResetColor();

                            ringBuffer.Busy.Release();
                            ringBuffer.Empty.Release();
                        }
                        else
                        {
                            ringBuffer.Full.Release();
                            ringBuffer.Busy.Release();
                        }
                        break;
                    case CharType.Other:
                        if (!char.IsLetterOrDigit(symbol))
                        {
                            ringBuffer.Tail = (ringBuffer.Tail + 1) % ringBuffer.Buffer.Length;
                            
                            Console.ForegroundColor = ConsoleColor.Blue;
                            Console.CursorLeft -= 1;
                            Console.Write(symbol);
                            Console.ResetColor();

                            ringBuffer.Busy.Release();
                            ringBuffer.Empty.Release();
                        }
                        else
                        {
                            ringBuffer.Full.Release();
                            ringBuffer.Busy.Release();
                        }
                        break;
                    default:
                        break;
                }
            }
        }
        static void Main(string[] args)
        {
            var buffer = new RingBuffer(10);

            var producer = new Task(() => ProducerAction(ref buffer));

            var consumer1 = new Task(() => ConsumerAction(ref buffer, CharType.Digit));

            var consumer2 = new Task(() => ConsumerAction(ref buffer, CharType.Letter));

            var consumer3 = new Task(() => ConsumerAction(ref buffer, CharType.Other));

            producer.Start();
            consumer1.Start();
            consumer2.Start();
            consumer3.Start();


            Task.WaitAll(producer, consumer1, consumer2, consumer3);
            //Console.OutputEncoding = Encoding.UTF8;
            //Console.Write("Вводите символы : ");



        }
    }
}
