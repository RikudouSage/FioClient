#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "common.h"

typedef struct {
	char* account_number;
	char* bank_code;
	char* currency;
	char* iban;
	char* bic;
} FioClientAccount;

typedef struct {
	FioClientAccount* items;
	size_t length;
} FioClientAccounts;

typedef struct {
	int64_t id;
	char* account_number;
	char* date;
	char* amount;
	char* currency;
	char* counterparty_account;
	char* counterparty_name;
	char* counterparty_bank_code;
	char* counterparty_bank_name;
	char* constant_symbol;
	char* variable_symbol;
	char* specific_symbol;
	char* user_identity;
	char* transaction_type;
	char* performed_by;
	char* additional_info;
	char* comment;
	char* bic;
	int64_t* instruction_id;
	char* payer_reference;
	bool local_only;
} FioClientTransaction;

typedef struct {
	FioClientTransaction* items;
	size_t length;
} FioClientTransactions;
