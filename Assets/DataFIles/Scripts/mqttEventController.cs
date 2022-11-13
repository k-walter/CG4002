using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class mqttEventController : MonoBehaviour
{
    public string nameController = "Event Controller";
    public string tagOfTheMQTTReceiver="MQTTEvent";
    public mqttEventReceiver _eventSender;

    public mqttSend mqttSend;

    // Events
    public Effects effects;
    public SoundManagerScript soundManagerScript;
    public OpponentActions opponentActions;
    public GrenadeThrower grenadeThrower;
    public GrenadeThrower opponentGrenadeThrower;
    public GameObject shieldedBorder;
    public GameObject opponentShield;

    void Start()
    {
        _eventSender=GameObject.FindGameObjectsWithTag(tagOfTheMQTTReceiver)[0].gameObject.GetComponent<mqttEventReceiver>();
        _eventSender.OnMessageArrived += OnMessageArrivedHandler;
    }

    private void OnMessageArrivedHandler(string newMsg)
    {
        var gameEvent = JsonUtility.FromJson<mqttEvent>(newMsg);
        if (gameEvent.player == 2 && gameEvent.action == "shield") {
            effects.SetShieldActive();
            soundManagerScript.PlayShieldSound();
        } else if (gameEvent.player == 2 && gameEvent.action == "shoot") {
            soundManagerScript.PlayFireSound();
        } else if (gameEvent.player == 2 && gameEvent.action == "shot") { // change back to damaged
            effects.SetDamagedTime();
            soundManagerScript.PlayPlayerHitSound();
        } else if (gameEvent.player == 2 && gameEvent.action == "shieldAvailable") {
            effects.SetShieldCooldownFill();
            shieldedBorder.SetActive(false);
        } else if (gameEvent.player == 1 && gameEvent.action == "shieldAvailable") {
            opponentShield.SetActive(false);
        } else if (gameEvent.player == 1 && gameEvent.action == "shot") { // change back to damaged
            effects.OpponentDamaged();
        } else if (gameEvent.player == 1 && gameEvent.action == "shield") {
            effects.SetOpponentShieldActive();
        } else if (gameEvent.player == 2 && gameEvent.action == "checkFov") {
            if (effects.checkOpponentVisible()) {
                Output output = new Output();
                output.player = 1;
                output.time = gameEvent.time;
                output.inFov = true;
                output.rnd = gameEvent.rnd;
                string json = JsonUtility.ToJson(output);
                mqttSend.setMessagePublish(json);
                mqttSend.Publish();
            } else {
                Output output = new Output();
                output.player = 1;
                output.time = gameEvent.time;
                output.inFov = false;
                output.rnd = gameEvent.rnd;
                string json = JsonUtility.ToJson(output);
                mqttSend.setMessagePublish(json);
                mqttSend.Publish();
            }
        } else if (gameEvent.player == 2 && gameEvent.action == "grenade") {
            grenadeThrower.ThrowGrenade();
        } else if (gameEvent.player == 2 && gameEvent.action == "grenaded") {
            opponentGrenadeThrower.ThrowGrenade();
            effects.OpponentDamaged();
        } else if (gameEvent.player == 2 && gameEvent.action == "done") {
            effects.setActionCount();
        }
    }
}

[System.Serializable]
public class mqttEvent 
{
    public int player;
    public ulong time;
    public string action;
    public int rnd;
}

[System.Serializable]
public class Output
{
    public int player;
    public ulong time;
    public bool inFov;
    public int rnd;
}