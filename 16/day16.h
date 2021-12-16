#ifndef DAY_16_H
#define DAY_16_H

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include "common/testing.h"
#include "common/vector.h"

#define WORD_SIZE 4
typedef unsigned short word_t;

enum PacketType {
	ADD = 0,
	MUL = 1,
	MIN = 2,
	MAX = 3,
	LOAD = 4,
	GT = 5,
	LT = 6,
	EQ = 7
};

struct packet_t;

TEMPLATE_VECTOR(struct packet_t) PacketVector;

typedef struct packet_t {
    word_t version;
    word_t type_id;
    unsigned long value;
    PacketVector subpackets;
} Packet;

typedef struct {
    size_t bitcount;
    char * line;
    char * line_pos;
    word_t buffer;
    word_t mask;
} BitReader;

BitReader GetBitReader(const char * line);
unsigned short InterpretByte(char c_repr);
int NextBit(BitReader * reader, word_t * bit);
void ClearBitReader(BitReader * reader);

unsigned int ReadInt(BitReader * br, size_t n_bits);

Packet ReadPacket(BitReader * reader);
void ReadLiteral(BitReader * reader, Packet * packet);
void ReadSubPackets(BitReader * reader, Packet * packet);

void PrintPacket(Packet * packet, size_t indentation);
void ClearPacket(Packet * packet);

size_t AccumulatePacketVersions(Packet * packet);

typedef void(*AccumulativeOperation)(long *, long);
typedef bool(*LogicalOperation)(long, long);

void Op_Sum(long * x, long y);
void Op_Mul(long * x, long y);
void Op_Min(long * x, long y);
void Op_Max(long * x, long y);
bool Op_Gt(long x, long y);
bool Op_Lt(long x, long y);
bool Op_Eq(long x, long y);

long int Operate(Packet * packet);

// Solving
typedef long int solution_t;
solution_t SolvePart1(const bool is_test);
solution_t SolvePart2(const bool is_test);

DECLARE_TEST(1);
DECLARE_TEST(2);

#endif
