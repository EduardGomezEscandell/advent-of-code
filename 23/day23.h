#ifndef DAY_23_H
#define DAY_23_H

#include "routing.h"
#include "gamestate.h"

#include "common/testing.h"
#include "common/vector.h"

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>

typedef cost_t solution_t;

// Solving
solution_t SolvePart1(const bool is_test);
solution_t SolvePart2(const bool is_test);

DECLARE_TEST(1);
DECLARE_TEST(2);

#endif
