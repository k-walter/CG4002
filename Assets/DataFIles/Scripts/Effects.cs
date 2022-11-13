using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using TMPro;

public class Effects : MonoBehaviour
{
    public float shieldCooldownTime = 10;
    public bool isShieldCooldown = false;
    public float shieldCooldownFill = 0;

    public float opponentShieldCooldownTime = 10;
    public bool isOpponentShieldCooldown = false;
    public float opponentShieldCooldownFill = 0;
    public int opponentShieldHealth = 0;

    public float damagedTime = 0;
    public GameObject damagedBorder;

    public bool isOpponentVisible;
    public OpponentActions opponentActions;

    public ShieldUI shieldUI;

    public int actionCount = 2;
    public float actionTime = 1;
    public GameObject actionBorder;

    // Use this for initialization
    void Start()
    {

    }

    // Update is called once per frame
    void Update()
    {
        Shield();
        OpponentShield();
        Damaged();
        Action();
    }

    public void SetShieldActive()
    {
        if (isShieldCooldown == true)
        {
           // can't use shield if it's on cooldown
        } else
        {
            isShieldCooldown = true;
            shieldCooldownFill = 1;
            shieldUI.SetShieldCooldownImage(shieldCooldownFill);
        }
    }

    public void Shield()
    {
        if (isShieldCooldown)
        {
            shieldCooldownFill -= 1 / shieldCooldownTime * Time.deltaTime;
            shieldUI.SetShieldCooldownImage(shieldCooldownFill);

            if (shieldCooldownFill <= 0)
            {
                shieldCooldownFill = 0;
                shieldUI.SetShieldCooldownImage(shieldCooldownFill);
                isShieldCooldown = false;
            }
        }
    }

    public void SetShieldCooldownFill() {
        shieldCooldownFill = 0;
    }

    public void SetOpponentShieldActive()
    {
        if (isOpponentShieldCooldown == true)
        {
            // can't use shield if it's on cooldown
        }
        else
        {
            isOpponentShieldCooldown = true;
            opponentShieldCooldownFill = 1;;
        }
    }

    public void OpponentShield()
    {
        if (isOpponentShieldCooldown)
        {
            opponentShieldCooldownFill -= 1 / opponentShieldCooldownTime * Time.deltaTime;

            if (opponentShieldCooldownFill <= 0)
            {
                opponentShieldCooldownFill = 0;
                isOpponentShieldCooldown = false;
            }
        }
    }

    public void SetDamagedTime() {
        damagedTime = 1;
    }

    public void Damaged()
    {
        if (damagedTime > 0) {
            damagedBorder.SetActive(true);
            damagedTime -= Time.deltaTime;
        } else {
            damagedBorder.SetActive(false);
        }
    }

    public void OpponentDamaged()
    {
        if (isOpponentVisible)
        {
            opponentActions.BloodSplash();
        }
    }

    public void SetOpponentVisible()
    {
        isOpponentVisible = true;
    }

    public void SetOpponentNotVisible()
    {
        isOpponentVisible = false;
    }

    public bool checkOpponentVisible() {
        if (isOpponentVisible) {
            return true;
        } else {
            return false;
        }
    }

    public void setActionCount() {
        actionCount = 0;
    }

    public void setActionTime() {
        actionTime = 1;
    }

    public void Action() {
        if (actionCount < 2) {
            if (actionTime > 0.5) {
                actionBorder.SetActive(true);
                actionTime -= Time.deltaTime;
            } else if (actionTime > 0) {
                actionBorder.SetActive(false);
                actionTime -= Time.deltaTime;
            } else {
                actionCount += 1;
                setActionTime();
            }
        }
    }
}