#! /usr/bin/env python
# coding: utf-8

import random
from optparse import OptionParser

# 进程切换策略
SCHED_SWITCH_ON_IO = 'SWITCH_ON_IO'         # 有IO进行切换
SCHED_SWITCH_ON_END = 'SWITCH_ON_END'       # 进程结束进行切换

# io完成后的策略
IO_RUN_LATER = 'IO_RUN_LATER'               # 稍后继续执行
IO_RUN_IMMEDIATE = 'IO_RUN_IMMEDIATE'       # 立即执行


# 进程状态
STATE_RUNNING = 'RUNNING'
STATE_READY = 'READY'
STATE_DONE = 'DONE'
STATE_WAIT = 'WAITING'

# 进程结构属性
PROC_CODE = 'code_'  # 程序指令
PROC_PC = 'pc_'
PROC_ID = 'pid_'
PROC_STATE = 'proc_state_'

# 进程行为
DO_COMPUTE = 'cpu'
DO_IO = 'io'


class scheduler:
    def __init__(self, process_switch_behavior, io_done_behavior, io_length):
        """_summary_

        Args:
            process_switch_behavior (int): 进程切换策略
            io_done_behavior (int): io完成后的策略
            io_length (int): io时间周期长度
        """
        self.proc_info = {}
        self.process_switch_behavior = process_switch_behavior
        self.io_done_behavior = io_done_behavior
        self.io_length = io_length
        return

    # 新建一个进程
    def new_process(self):
        proc_id = len(self.proc_info)
        self.proc_info[proc_id] = {}
        self.proc_info[proc_id][PROC_PC] = 0
        self.proc_info[proc_id][PROC_ID] = proc_id
        self.proc_info[proc_id][PROC_CODE] = []
        self.proc_info[proc_id][PROC_STATE] = STATE_READY
        return proc_id

    # 加载程序文件，分配进程，加载指令
    def load_file(self, progfile):
        fd = open(progfile)
        proc_id = self.new_process()

        for line in fd:
            tmp = line.split()
            if len(tmp) == 0:
                continue
            opcode = tmp[0]
            if opcode == 'compute':
                assert(len(tmp) == 2)
                for i in range(int(tmp[1])):
                    self.proc_info[proc_id][PROC_CODE].append(DO_COMPUTE)
            elif opcode == 'io':
                assert(len(tmp) == 1)
                self.proc_info[proc_id][PROC_CODE].append(DO_IO)
        fd.close()
        return

    # 加载程序描述，给出全部指令和cpu指令的占比
    def load(self, program_description):
        proc_id = self.new_process()
        tmp = program_description.split(':')
        if len(tmp) != 2:
            print('Bad description (%s): Must be number <x:y>' %
                  program_description)
            print('  where X is the number of instructions')
            print('  and Y is the percent change that an instruction is CPU not IO')
            exit(1)

        num_instructions, chance_cpu = int(tmp[0]), float(tmp[1])/100.0
        for i in range(num_instructions):
            if random.random() < chance_cpu:
                self.proc_info[proc_id][PROC_CODE].append(DO_COMPUTE)
            else:
                self.proc_info[proc_id][PROC_CODE].append(DO_IO)
        return

    # 使进程就绪
    def move_to_ready(self, expected, pid=-1):
        if pid == -1:
            pid = self.curr_proc
        assert(self.proc_info[pid][PROC_STATE] == expected)
        self.proc_info[pid][PROC_STATE] = STATE_READY
        return

    # 使进程进入等待状态，必须是运行状态
    def move_to_wait(self, expected):
        assert(self.proc_info[self.curr_proc][PROC_STATE] == expected)
        self.proc_info[self.curr_proc][PROC_STATE] = STATE_WAIT
        return

    # 使进程run起来，必须是就绪状态
    def move_to_running(self, expected):
        assert(self.proc_info[self.curr_proc][PROC_STATE] == expected)
        self.proc_info[self.curr_proc][PROC_STATE] = STATE_RUNNING
        return

    # 使进程停下来，必须是运行状态
    def move_to_done(self, expected):
        assert(self.proc_info[self.curr_proc][PROC_STATE] == expected)
        self.proc_info[self.curr_proc][PROC_STATE] = STATE_DONE
        return

    # 执行下一条指令
    def next_proc(self, pid=-1):
        if pid != -1:
            self.curr_proc = pid
            self.move_to_running(STATE_READY)
            return
        for pid in range(self.curr_proc + 1, len(self.proc_info)):
            if self.proc_info[pid][PROC_STATE] == STATE_READY:
                self.curr_proc = pid
                self.move_to_running(STATE_READY)
                return
        for pid in range(0, self.curr_proc + 1):
            if self.proc_info[pid][PROC_STATE] == STATE_READY:
                self.curr_proc = pid
                self.move_to_running(STATE_READY)
                return
        return

    # 获取进程数量
    def get_num_processes(self):
        return len(self.proc_info)

    # 获取进程的指令数
    def get_num_instructions(self, pid):
        return len(self.proc_info[pid][PROC_CODE])

    # 获取指定指令
    def get_instruction(self, pid, index):
        return self.proc_info[pid][PROC_CODE][index]

    # 获取活跃的进程数量
    def get_num_active(self):
        num_active = 0
        for pid in range(len(self.proc_info)):
            if self.proc_info[pid][PROC_STATE] != STATE_DONE:
                num_active += 1
        return num_active

    # 获取可运行的次数
    def get_num_runnable(self):
        num_active = 0
        for pid in range(len(self.proc_info)):
            if self.proc_info[pid][PROC_STATE] == STATE_READY or \
                   self.proc_info[pid][PROC_STATE] == STATE_RUNNING:
                num_active += 1
        return num_active
    
    # 获取当前正在执行io的个数
    def get_ios_in_flight(self, current_time):
        num_in_flight = 0
        for pid in range(len(self.proc_info)):
            for t in self.io_finish_times[pid]:
                if t > current_time:
                    num_in_flight += 1
        return num_in_flight

    def check_for_switch(self):
        return

    # 打印空格
    def space(self, num_columns):
        for i in range(num_columns):
            print('%10s' % ' '),

    # 校验当前指令是否还有指令需要执行
    def check_if_done(self):
        if len(self.proc_info[self.curr_proc][PROC_CODE]) == 0:
            if self.proc_info[self.curr_proc][PROC_STATE] == STATE_RUNNING:
                self.move_to_done(STATE_RUNNING)
                self.next_proc()
        return

    # 运行
    def run(self):
        # 时钟周期
        clock_tick = 0

        if len(self.proc_info) == 0:
            return

        # 跟踪每个进程IO完成的次数
        self.io_finish_times = {}
        for pid in range(len(self.proc_info)):
            self.io_finish_times[pid] = []

        # 让第一个进程处于活跃状态
        self.curr_proc = 0
        self.move_to_running(STATE_READY)

        # 输出：每列标题
        print('%s' % 'Time'),
        for pid in range(len(self.proc_info)):
            print('%10s' % ('PID:%2d' % (pid))),
        print('%10s' % 'CPU'),
        print('%10s' % 'IOs'),
        print('')

        # init statistics
        io_busy = 0
        cpu_busy = 0

        while self.get_num_active() > 0:
            clock_tick += 1

            # 检查io，将io完成的进程改到就绪状态
            io_done = False
            for pid in range(len(self.proc_info)):
                if clock_tick in self.io_finish_times[pid]:
                    io_done = True
                    self.move_to_ready(STATE_WAIT, pid)
                    
                    if self.io_done_behavior == IO_RUN_IMMEDIATE:
                        # IO立即执行策略，将进程切换到IO完成的进程
                        if self.curr_proc != pid:
                            if self.proc_info[self.curr_proc][PROC_STATE] == STATE_RUNNING:
                                self.move_to_ready(STATE_RUNNING)
                        self.next_proc(pid)
                    else:
                        if self.process_switch_behavior == SCHED_SWITCH_ON_END and self.get_num_runnable() > 1:
                            # 等当前进程IO结束后，修改当前进程为运行状态
                            self.next_proc(pid)
                        if self.get_num_runnable() == 1:
                            # 运行这个唯一可运行的进程
                            self.next_proc(pid)
                    # 校验当前进程是否还有指令，没有则切换到下一个进程
                    self.check_if_done()

            # 如果进程活跃，并且还有指令，拿到当前指令
            instruction_to_execute = ''
            if self.proc_info[self.curr_proc][PROC_STATE] == STATE_RUNNING and \
                    len(self.proc_info[self.curr_proc][PROC_CODE]) > 0:
                instruction_to_execute = self.proc_info[self.curr_proc][PROC_CODE].pop(0)
                cpu_busy += 1

            # 输出指令周期数，如果由IOdone，则加*
            if io_done:
                print('%3d*' % clock_tick),
            else:
                print('%3d ' % clock_tick),
            # 输出当前时间周期下，各进程运行的指令情况
            for pid in range(len(self.proc_info)):
                if pid == self.curr_proc and instruction_to_execute != '':
                    print('%10s' % ('RUN:'+instruction_to_execute)),
                else:
                    print('%10s' % (self.proc_info[pid][PROC_STATE])),
            # 输出当前时间周期CPU是否执行
            if instruction_to_execute == '':
                print('%10s' % ' '),
            else:
                print('%10s' % 1),
            
            # 输出当前正在执行的io个数
            num_outstanding = self.get_ios_in_flight(clock_tick)
            if num_outstanding > 0:
                print('%10s' % str(num_outstanding)),
                io_busy += 1
            else:
                print('%10s' % ' '),
            print('')

            # 对于io指令，切换到等待状态。
            if instruction_to_execute == DO_IO:
                self.move_to_wait(STATE_RUNNING)
                self.io_finish_times[self.curr_proc].append(
                    clock_tick + self.io_length)
                # 这种情况切换进程
                if self.process_switch_behavior == SCHED_SWITCH_ON_IO:
                    self.next_proc()

            # 检查当前运行的东西是否超出指令
            self.check_if_done()
        return (cpu_busy, io_busy, clock_tick)

#
# PARSE ARGUMENTS
#

parser = OptionParser()
parser.add_option('-s', '--seed', default=0, help='随机种子', action='store', type='int', dest='seed')
parser.add_option('-l', '--processlist', default='',
                  help='以逗号分隔的要运行的进程列表，格式为X1：Y1，X2：Y2，...，其中X是该进程应运行的指令数，Y是该指令将运行的概率（从0到100）,指令包括使用CPU或进行IO',
                  action='store', type='string', dest='process_list')
parser.add_option('-L', '--iolength', default=5, help='IO花费时间', action='store', type='int', dest='io_length')
parser.add_option('-S', '--switch', default='SWITCH_ON_IO',
                  help='当进程发出IO时,系统的反应D',
                  action='store', type='string', dest='process_switch_behavior')
parser.add_option('-I', '--iodone', default='IO_RUN_LATER',
                  help='IO结束时的行为类型: IO_RUN_LATER:自然切换到这个进程(例如:取决于进程切换行为), IO_RUN_IMMEDIATE：立即切换到这个进程',
                  action='store', type='string', dest='io_done_behavior')
parser.add_option('-c', help='计算结果', action='store_true', default=False, dest='solve')
parser.add_option('-p', '--printstats', help='打印统计数据； 仅与-c参数一起使用是有效', action='store_true', default=False, dest='print_stats')
(options, args) = parser.parse_args()

random.seed(options.seed)

assert(options.process_switch_behavior == SCHED_SWITCH_ON_IO or \
       options.process_switch_behavior == SCHED_SWITCH_ON_END)
assert(options.io_done_behavior == IO_RUN_IMMEDIATE or \
       options.io_done_behavior == IO_RUN_LATER)

s = scheduler(options.process_switch_behavior, options.io_done_behavior, options.io_length)

# example process description (10:100,10:100)
for p in options.process_list.split(','):
    s.load(p)

if options.solve == False:
    print('进程执行信息信息:')
    for pid in range(s.get_num_processes()):
        print('进程 %d' % pid)
        for inst in range(s.get_num_instructions(pid)):
            print('  %s' % s.get_instruction(pid, inst))
        print('')
    print('重要行为:')
    print('  系统切换时机'),
    if options.process_switch_behavior == SCHED_SWITCH_ON_IO:
        print('当前进程结束或发起IO操作')
    else:
        print('当前进程结束')
    print('  io结束后，发起IO的进程行为'),
    if options.io_done_behavior == IO_RUN_IMMEDIATE:
        print('立即执行')
    else:
        print('等待调度')
    print('')
    exit(0)

(cpu_busy, io_busy, clock_tick) = s.run()

if options.print_stats:
    print('')
    print('Stats: Total Time %d' % clock_tick)
    print('Stats: CPU Busy %d (%.2f%%)' % (cpu_busy, 100.0 * float(cpu_busy)/clock_tick))
    print('Stats: IO Busy  %d (%.2f%%)' % (io_busy, 100.0 * float(io_busy)/clock_tick))
    print('')