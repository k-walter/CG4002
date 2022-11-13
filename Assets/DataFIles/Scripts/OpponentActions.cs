using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class OpponentActions : MonoBehaviour
{

    public GameObject bloodSplash;

    // Start is called before the first frame update
    void Start()
    {
        
    }

    // Update is called once per frame
    void Update()
    {

    }

    public void BloodSplash()
    {
        Instantiate(bloodSplash, transform.position, Quaternion.identity);
    }
}
