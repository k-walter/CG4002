using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class ShieldUI : MonoBehaviour
{
    //Shield bar
    public Image[] shieldPoints;

    public void SetShield(int shield)
    {
        for (int i = 0; i < shieldPoints.Length; i++)
        {
            shieldPoints[i].enabled = !DisplayShieldPoint(shield, i);
        }
    }

    bool DisplayShieldPoint(int shield, int pointNumber)
    {
        return pointNumber >= shield;
    }

    // Shield cooldown
    public Image shieldImage;

    public void SetShieldCooldownImage(float fill)
    {
        shieldImage.fillAmount = fill;
    }
}
