using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class GrenadeBar : MonoBehaviour
{
    public Image[] grenadePoints;

    public void SetGrenade(int grenade)
    {
        for (int i = 0; i < grenadePoints.Length; i++)
        {
            grenadePoints[i].enabled = !DisplayGrenadePoint(grenade, i);
        }
    }

    bool DisplayGrenadePoint(int grenade, int pointNumber)
    {
        return pointNumber >= grenade;
    }
}
