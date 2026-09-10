#pragma once

#include <stdlib.h>
#include <string.h>

static __thread char* fioclient_last_error;

static void fioclient_clear_last_error(void) {
	if (fioclient_last_error) {
		free(fioclient_last_error);
		fioclient_last_error = NULL;
	}
}

static void fioclient_set_last_error_copy(const char* message) {
	fioclient_clear_last_error();
	if (!message) return;

	size_t length = strlen(message);
	char* copy = (char*)malloc(length + 1);
	if (!copy) return;
	memcpy(copy, message, length + 1);
	fioclient_last_error = copy;
}

static size_t fioclient_get_last_error(char* buffer, size_t buffer_length) {
	if (!fioclient_last_error) {
		if (buffer && buffer_length) buffer[0] = '\0';
		return 1;
	}

	size_t required = strlen(fioclient_last_error) + 1;
	if (buffer && buffer_length) {
		size_t count = required <= buffer_length ? required - 1 : buffer_length - 1;
		memcpy(buffer, fioclient_last_error, count);
		buffer[count] = '\0';
	}
	return required;
}
