using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class AmmoBar : MonoBehaviour
{
    public Image[] ammoPoints;

    public void SetAmmo(int ammo)
    {
        for (int i = 0; i < ammoPoints.Length; i++)
        {
            ammoPoints[i].enabled = !DisplayAmmoPoint(ammo, i);
        }
    }

    bool DisplayAmmoPoint(int ammo, int pointNumber)
    {
        return pointNumber >= ammo;
    }
}
