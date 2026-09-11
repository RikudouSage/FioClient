#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libfioclient.h"

static int failures;

#define EXPECT(condition) do { \
	if (!(condition)) { \
		fprintf(stderr, "%s:%d: check failed: %s\n", __FILE__, __LINE__, #condition); \
		failures++; \
	} \
} while (0)

static void expect_error_contains(const char* expected) {
	size_t required = FioClientGetLastError(NULL, 0);
	EXPECT(required > 1);

	char* message = malloc(required);
	EXPECT(message != NULL);
	if (message == NULL) {
		return;
	}

	EXPECT(FioClientGetLastError(message, required) == required);
	EXPECT(strstr(message, expected) != NULL);

	free(message);
}

static void test_client_validation(void) {
	FioClientOptions options = {0};

	EXPECT(FioClientNew(NULL, options) == FioClientFailure);
	expect_error_contains("out is NULL");

	FioClientHandle client = 123;
	EXPECT(FioClientNew(&client, options) == FioClientFailure);
	EXPECT(client == 0);
	expect_error_contains("options.database_path is NULL");

	options.database_path = "test.sqlite3";
	EXPECT(FioClientNew(&client, options) == FioClientFailure);
	EXPECT(client == 0);
	expect_error_contains("options.database_key is NULL");
}

static void test_context_and_handles(void) {
	EXPECT(FioClientNewContext(NULL) == FioClientFailure);
	expect_error_contains("out is NULL");

	FioClientContextHandle context = 0;
	EXPECT(FioClientNewContext(&context) == FioClientSuccess);
	EXPECT(context != 0);
	EXPECT(FioClientGetLastError(NULL, 0) == 1);

	FioClientAccounts accounts = {
		.items = (FioClientAccount*)1,
		.length = 99,
	};
	EXPECT(FioClientGetAccounts(999999, context, &accounts) == FioClientFailure);
	EXPECT(accounts.items == NULL);
	EXPECT(accounts.length == 0);
	expect_error_contains("handle 999999 not found");

	EXPECT(FioClientCloseHandle(context) == FioClientSuccess);
	EXPECT(FioClientCloseHandle(context) == FioClientFailure);
	expect_error_contains("not registered");
}

static void test_output_validation_and_free(void) {
	EXPECT(FioClientGetAccounts(0, 0, NULL) == FioClientFailure);
	expect_error_contains("out is NULL");

	EXPECT(FioClientOpenAccount(0, 0, NULL, NULL) == FioClientFailure);
	expect_error_contains("out is NULL");

	EXPECT(FioClientGetTransactions(0, 0, NULL) == FioClientFailure);
	expect_error_contains("out is NULL");
	EXPECT(FioClientLoadTransactionsByDate(0, 0, NULL, NULL, NULL) == FioClientFailure);
	expect_error_contains("out is NULL");

	FioClientAccount account = {0};
	FioClientAccounts accounts = {0};
	FioClientTransaction transaction = {0};
	FioClientTransactions transactions = {0};

	FioClientFreeAccount(NULL);
	FioClientFreeAccount(&account);
	FioClientFreeAccounts(NULL);
	FioClientFreeAccounts(&accounts);
	FioClientFreeTransaction(NULL);
	FioClientFreeTransaction(&transaction);
	FioClientFreeTransactions(NULL);
	FioClientFreeTransactions(&transactions);
}

int main(void) {
	test_client_validation();
	test_context_and_handles();
	test_output_validation_and_free();

	if (failures != 0) {
		fprintf(stderr, "%d C API test(s) failed\n", failures);
		return EXIT_FAILURE;
	}

	puts("C API tests passed");

	return EXIT_SUCCESS;
}
