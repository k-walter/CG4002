#ifndef CNN_MODEL
#define CNN_MODEL

#include "cnn.h"
#include "dense_layer.h"

typedef int FUNCTION_SELECT_BIT_WIDTH;

void cnn_action_detection(
		FUNCTION_SELECT_BIT_WIDTH function_select,
		int user_number,
		float data[INPUT_DEPTH],
		float raw_output[DENSE_OUTPUT_NODES],
		float weights_and_bias[CNN_KERNEL_COUNT * CNN_KERNEL_LENGTH * CNN_KERNEL_DEPTH]);

#endif
