using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using TMPro;

public class mqttStateController : MonoBehaviour
{
    public string nameController = "State Controller";
    public string tagOfTheMQTTReceiver="MQTTState";
    public mqttStateReceiver _eventSender;

    // Standard UI
    public int maxHealth = 100;
    public int maxAmmo = 6;
    public int maxGrenade = 2;
    public int maxShield = 3;
    public int opponentMaxHealth = 100;

    public GameObject shieldedBorder;
    public GameObject opponentShield;

    public HealthBar healthBar;
    public AmmoBar ammoBar;
    public GrenadeBar grenadeBar;
    public ShieldUI shieldUI;
    public OpponentHealthBar opponentHealthBar;
    public TMP_Text player1Score;
    public TMP_Text player2Score;

    public GameObject skullNormal;
    public GameObject skullHit;

    void Start()
    {
        _eventSender=GameObject.FindGameObjectsWithTag(tagOfTheMQTTReceiver)[0].gameObject.GetComponent<mqttStateReceiver>();
        _eventSender.OnMessageArrived += OnMessageArrivedHandler;

        healthBar.SetHealth(maxHealth);
        ammoBar.SetAmmo(maxAmmo);
        grenadeBar.SetGrenade(maxGrenade);
        shieldUI.SetShield(maxShield);
        opponentHealthBar.SetOpponentHealth(opponentMaxHealth);
        player1Score.text = "0";
        player2Score.text = "0";

    }

    private void OnMessageArrivedHandler(string newMsg)
    {
        var gameState = JsonUtility.FromJson<State>(newMsg);
        healthBar.SetHealth(gameState.p2.hp);
        ammoBar.SetAmmo(gameState.p2.bullets);
        grenadeBar.SetGrenade(gameState.p2.grenades);
        shieldUI.SetShield(gameState.p2.num_shield);
        opponentHealthBar.SetOpponentHealth(gameState.p1.hp);
        player1Score.text = gameState.p1.num_deaths.ToString();
        player2Score.text = gameState.p2.num_deaths.ToString();
        if (gameState.p2.shield_health > 0) {
            shieldedBorder.SetActive(true);
        } else {
            shieldedBorder.SetActive(false);
        }
        if (gameState.p1.shield_health > 0) {
            opponentShield.SetActive(true);
        } else {
            opponentShield.SetActive(false);
        }
        if (gameState.p1.hp > 50) {
            skullNormal.SetActive(true);
            skullHit.SetActive(false);
        } else {
            skullHit.SetActive(true);
            skullNormal.SetActive(false);
        }
    }
}

[System.Serializable]
public class State 
{
    public player p1;
    public player p2;
}

[System.Serializable]
public class player
{
    public int hp;
    public string action;
    public int bullets;
    public int grenades;
    public double shield_time;
    public int shield_health;
    public int num_deaths;
    public int num_shield;
}