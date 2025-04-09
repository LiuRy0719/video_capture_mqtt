#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mosquitto.h"
#define MQTT_BROKER "112.6.224.25"
#define MQTT_PORT 20042
#define MQTT_TOPIC "testtopic"
#define Qos 1
void on_connect(struct mosquitto *mosq, void *obj, int rc)
{
    if (rc != MOSQ_ERR_SUCCESS)
    {
        fprintf(stderr, "Failed to connect to MQTT broker with reason code:%d\n", rc);
        exit(EXIT_FAILURE);
    }
    fprintf(stdout, "Connected to MQTT broker!\n");
    mosquitto_subscribe(mosq, NULL, MQTT_TOPIC, Qos);
}
int main()
{
    struct mosquitto *mosq = NULL;
    mosquitto_lib_init();
    mosq = mosquitto_new(NULL, true, NULL);
    mosquitto_connect_callback_set(mosq, on_connect);
    if (mosquitto_connect(mosq, MQTT_BROKER, MQTT_PORT, 60) != MOSQ_ERR_SUCCESS)
    {
        fprintf(stderr, "Failed to connect to MQTT broker!\n");
        exit(EXIT_FAILURE);
    }
    while (1)
    {
        int mid = 0;
        char *mes = "hello Lry";
        mosquitto_publish(mosq, &mid, MQTT_TOPIC, strlen(mes), mes, Qos, false);
        // printf("id=%d\n", mid);
        fprintf(stdout,"id=%d\n",mid);
        sleep(2);
    }
    mosquitto_disconnect(mosq);
    mosquitto_destroy(mosq);
    mosquitto_lib_cleanup();
}
