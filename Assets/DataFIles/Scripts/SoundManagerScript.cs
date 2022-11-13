using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class SoundManagerScript : MonoBehaviour
{

    public AudioClip fireSound, playerHitSound, shieldSound, grenadeSound;
    public AudioSource audioSrc;

    // Start is called before the first frame update
    void Start()
    {

    }

    // Update is called once per frame
    void Update()
    {
        
    }

    public void PlayFireSound()
    {
        audioSrc.PlayOneShot(fireSound);
    }

    public void PlayPlayerHitSound()
    {
        audioSrc.PlayOneShot(playerHitSound);
    }

    public void PlayShieldSound()
    {
        audioSrc.PlayOneShot(shieldSound);
    }

    public void PlayGrenadeSound()
    {
        audioSrc.PlayOneShot(grenadeSound);
    }
}
