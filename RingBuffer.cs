using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace ParallelProgramming_lab1
{
    class RingBuffer
    {
        public char[] Buffer { get; set; }

        public int Head { get; set; }

        public int Tail { get; set; }

        public Semaphore Full { get; set; }

        public Semaphore Empty { get; set; }

        public Semaphore Busy { get; set; }

        public RingBuffer(int capacity)
        {
            Buffer = new char[capacity];

            Head = 0;

            Tail = 0;

            Full = new Semaphore(0, capacity);

            Empty = new Semaphore(capacity, capacity);

            Busy = new Semaphore(1, 1);
        }
    }
}
