#include <stdio.h>
#include <vector>
#include <iostream>
#include "cnn_model.h"

#include "cnn_model_test_dataset.h"
// #include "hls_math.h"

int confusion[DENSE_OUTPUT_NODES][DENSE_OUTPUT_NODES];
int accurate_count = 0;

int result;
float data[INPUT_DEPTH];
float raw_outputs[DENSE_OUTPUT_NODES];
float weights_and_bias[CNN_KERNEL_COUNT * CNN_KERNEL_LENGTH * CNN_KERNEL_DEPTH];

// void softmax(float input[], float output[]) {
// 	float sum = 0;
// 	for (int i = 0; i < DENSE_OUTPUT_NODES; i++) {
// 		sum += hls::exp(input[i]);
// 	}
// 	for (int i = 0; i < DENSE_OUTPUT_NODES; i++) {
// 		output[i] = hls::exp(input[i]) / sum;
// 	}
// }

void motionDetect(int user) {
	for (int data_index = 0; data_index < DATASET_SIZE; data_index++) {
		// feed the whole window
		for (int i = 0; i < INPUT_LENGTH; i++) {
			for (int j = 0; j < INPUT_DEPTH; j++) {
				data[j] = test_x[data_index][i][j];
			}
			cnn_action_detection(0,
				0,
				data,
				raw_outputs,
				weights_and_bias
			);
			result = argmax(raw_outputs);

			// // print raw outputs
			// printf("current raw outputs: \n");
			// for (int i = 0; i < DENSE_OUTPUT_NODES; i++) {
			// 	printf("%f ", raw_outputs[i]);
			// }
			// printf("\n");
		}

		// print raw outputs
		printf("current raw outputs: ");
		for (int i = 0; i < DENSE_OUTPUT_NODES; i++) {
			printf("%f\t", raw_outputs[i]);
		}
		printf("\n");

		// // compute softmax
		// float softmax_outputs[DENSE_OUTPUT_NODES];
		// softmax(raw_outputs, softmax_outputs);

		int GOLD_result;
		GOLD_result = argmax(test_y[data_index]);

		// add to confusion matrix
		confusion[GOLD_result][result]++;

		if (result == GOLD_result) {
			accurate_count += 1;
		}

		// printf("Data %d: predicted = %d vs GOLD = %d\n", data_index, result, GOLD_result);
		

		// reset
		cnn_action_detection(1, 0, data, raw_outputs, weights_and_bias);
	}

	// print accuracy
	printf("Accuracy: %f\n", float(accurate_count) / DATASET_SIZE);

	// print confusion matrix
	for (int i = 0; i < DENSE_OUTPUT_NODES; i++) {
		for (int j = 0; j < DENSE_OUTPUT_NODES; j++) {
			printf("%d\t", confusion[i][j]);
		}
		printf("\n");
	}
}

// void cnn_layer_test() {
//     CNN_DTYPE sample_cnn_output_buffer[CNN_OUTPUT_LENGTH][CNN_OUTPUT_DEPTH];
//     CNN_DTYPE sample_averaged[CNN_OUTPUT_DEPTH];

//     for (int i = 0; i < CNN_OUTPUT_LENGTH; i++) {
//         sample_cnn_output_buffer[CNN_OUTPUT_LENGTH-1][0] = 73.0;
//         compute_global_average_pool(sample_cnn_output_buffer, sample_averaged);
//         printf("sample_averaged[0] = %f\n", sample_averaged[0]);
//     }
// }

int main() {
	// cnn_layer_test();
	motionDetect(0);
}
