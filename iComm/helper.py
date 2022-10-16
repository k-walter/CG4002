from constants import MAX_16_BIT_SIGNED, MAX_16_BIT_UNSIGNED

# Helper function to unpack IMU sensor data
def unpack_glove_data_into_dict(glove_data):
    glove_dict = {
        "index": glove_data[0],
        "roll": bytes_to_uint16_t(glove_data[1:3]),
        "pitch": bytes_to_uint16_t(glove_data[3:5]),
        "yaw": bytes_to_uint16_t(glove_data[5:7]),
        "x": bytes_to_uint16_t(glove_data[7:9]),
        "y": bytes_to_uint16_t(glove_data[9:11]),
        "z": bytes_to_uint16_t(glove_data[11:13]),
    }
    return glove_dict

def bytes_to_uint16_t(bytes):
    val = (bytes[0] << 8) + bytes[1]
    if val > MAX_16_BIT_SIGNED:
        val = val - MAX_16_BIT_UNSIGNED
    return val