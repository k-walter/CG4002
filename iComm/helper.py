from constants import MAX_16_BIT_SIGNED, MAX_16_BIT_UNSIGNED

# Helper function to unpack IMU sensor data
def unpack_glove_data_into_dict(glove_data):
    glove_dict = {
        "index": bytes_to_uint16_t(glove_data[:2]),
        "roll": bytes_to_uint16_t(glove_data[2:4]),
        "pitch": bytes_to_uint16_t(glove_data[4:6]),
        "yaw": bytes_to_uint16_t(glove_data[6:8]),
        "x": bytes_to_uint16_t(glove_data[8:10]),
        "y": bytes_to_uint16_t(glove_data[10:12]),
        "z": bytes_to_uint16_t(glove_data[12:14]),
    }
    return glove_dict

def bytes_to_uint16_t(bytes):
    val = (bytes[0] << 8) + bytes[1]
    if val > MAX_16_BIT_SIGNED:
        val = val - MAX_16_BIT_UNSIGNED
    return val