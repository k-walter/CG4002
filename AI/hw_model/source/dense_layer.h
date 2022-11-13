/*
 * Dense Layers
 * Created: 4 Sep 2022
 * Author: Yue Junfeng
 */

#ifndef DENSE_LAYER
#define DENSE_LAYER

#include "cnn.h"

#define DENSE_INPUT_NODES CNN_OUTPUT_DEPTH
#define DENSE_OUTPUT_NODES 5

typedef float DENSE_DTYPE;

void compute_dense(
		CNN_DTYPE input_buffer[CNN_OUTPUT_DEPTH],
		DENSE_DTYPE dense_weights[DENSE_OUTPUT_NODES][DENSE_INPUT_NODES],
		DENSE_DTYPE dense_bias[DENSE_OUTPUT_NODES],
		DENSE_DTYPE dense_output[DENSE_OUTPUT_NODES]
		);

template<typename T>
T argmax(T input[]){
	DENSE_DTYPE max = input[0];
	int max_index = 0;
	ARGMAX: for (int i = 1; i < DENSE_OUTPUT_NODES; i++) {
		if (input[i] > max) {
			max = input[i];
			max_index = i;
		}
	}
	return max_index;
}

#endif
