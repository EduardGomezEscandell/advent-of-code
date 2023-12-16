#include "day16.hpp"
#include <doctest/doctest.h>

TEST_CASE("Day 16") {

  SUBCASE("Part 1, example") {
    Day16 solution{};
    solution.set_input("./data/16/example.txt");
    solution.load();
    REQUIRE_THROWS(solution.part1());
  }

  SUBCASE("Part 2, example") {
    Day16 solution{};
    solution.set_input("./data/16/example.txt");
    solution.load();
    REQUIRE_THROWS(solution.part2());
  }

  SUBCASE("Real data") {
    Day16 solution{};
    solution.set_input("./data/16/input.txt");
    solution.load();
    REQUIRE_THROWS(solution.part1());
    REQUIRE_THROWS(solution.part2());
  }
}