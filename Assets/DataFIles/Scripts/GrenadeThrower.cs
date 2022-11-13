using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class GrenadeThrower : MonoBehaviour
{

    public float throwForce = 20f;
    public GameObject grenadePrefab;
    
    public SoundManagerScript soundManagerScript;
    float countdown = 2;
    bool playGrenadeSound = false;

    // Start is called before the first frame update
    void Start()
    {
        
    }

    // Update is called once per frame
    void Update()
    {
        if (playGrenadeSound == true) {
            countdown -= Time.deltaTime;
            if (countdown <= 0f) {
                soundManagerScript.PlayGrenadeSound();
                playGrenadeSound = false;
                countdown = 2;
            }
        }
    }

    public void ThrowGrenade()
    {
        GameObject grenade = Instantiate(grenadePrefab, transform.position, transform.rotation);
        Rigidbody rb = grenade.GetComponent<Rigidbody>();
        rb.AddForce(transform.forward * throwForce, ForceMode.VelocityChange);
        playGrenadeSound = true;
    }
}
