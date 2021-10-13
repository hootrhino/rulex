// crt_assert.c
// compile by using: cl /W4 crt_assert.c
#include <stdio.h>
#include <assert.h>
#include <string.h>

void analyze_string(char *string); // Prototype

int main(void)
{
    char test1[] = "abc", *test2 = NULL, test3[] = "";

    printf("Analyzing string '%s'\n", test1);
    fflush(stdout);
    analyze_string(test1);
    printf("Analyzing string '%s'\n", test2);
    fflush(stdout);
    analyze_string(test2);
    printf("Analyzing string '%s'\n", test3);
    fflush(stdout);
    analyze_string(test3);
}

// Tests a string to see if it is NULL,
// empty, or longer than 0 characters.
void analyze_string(char *string)
{
    assert(string != NULL);     // Cannot be NULL
    assert(*string != '\0');    // Cannot be empty
    assert(strlen(string) > 2); // Length must exceed 2
}