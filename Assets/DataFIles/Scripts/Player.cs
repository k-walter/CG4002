// using System.Collections;
// using System.Collections.Generic;
// using UnityEngine;
// using UnityEngine.UI;
// using TMPro;

// public class Player : MonoBehaviour
// {

//     public int maxHealth = 100;
//     public int currentHealth;

//     public int maxAmmo = 6;
//     public int currentAmmo;

//     public int maxGrenade = 2;
//     public int currentGrenade;

//     public int maxShield = 3;
//     public int currentShield;
//     public float shieldCooldownTime = 10;
//     public bool isShieldCooldown = false;
//     public float shieldCooldownFill;
//     public GameObject shieldedBorder;
//     public Button shieldButton;
//     public int shieldHealth = 0;

//     public float opponentShieldCooldownTime = 10;
//     public bool isOpponentShieldCooldown = false;
//     public float opponentShieldCooldownFill = 0;
//     public int opponentShieldHealth = 0;
//     public GameObject opponentShield;

//     public int maxOpponentHealth = 100;
//     public int currentOpponentHealth;

//     public bool isDamaged = false;
//     public float damagedTime = 1;
//     public GameObject damagedBorder;

//     public int player1 = 0;
//     public TMP_Text player1Score;
//     public int player2 = 0;
//     public TMP_Text player2Score;

//     public bool isOpponentVisible;

//     public HealthBar healthBar;
//     public AmmoBar ammoBar;
//     public GrenadeBar grenadeBar;
//     public ShieldUI shieldUI;
//     public OpponentHealthBar opponentHealthBar;
//     public SoundManagerScript soundManagerScript;
//     public OpponentActions opponentActions;

//     // Use this for initialization
//     void Start()
//     {
//         currentHealth = maxHealth;
//         healthBar.SetMaxHealth(maxHealth);

//         currentAmmo = maxAmmo;
//         ammoBar.SetAmmo(maxAmmo);

//         currentGrenade = maxGrenade;
//         grenadeBar.SetGrenade(maxGrenade);

//         currentShield = maxShield;
//         shieldUI.SetShield(maxShield);

//         currentOpponentHealth = maxOpponentHealth;
//         opponentHealthBar.SetOpponentHealth(maxOpponentHealth);
//     }

//     // Update is called once per frame
//     void Update()
//     {
//         Health();
//         Ammo(); Use button to call this function instead
//         Grenade(); Use button to call this function instead
//         Shield();
//         OpponentShield();
//         Damaged();
//         Score();
//     }

//     void Health()
//     {
//         if (currentHealth <= 0)
//         {
//             soundManagerScript.PlayDeathSound();
//             currentHealth = maxHealth;
//             healthBar.SetHealth(currentHealth);
//             player1 += 1;
//         }
//         if (currentOpponentHealth <= 0)
//         {
//             soundManagerScript.PlayDeathSound();
//             currentOpponentHealth = maxOpponentHealth;
//             opponentHealthBar.SetOpponentHealth(currentOpponentHealth);
//             player2 += 1;
//         }
//     }

//     public void Ammo()
//     {
//         currentAmmo -= 1;
//         ammoBar.SetAmmo(currentAmmo);
//         if (opponentShieldHealth != 0)
//         {
//             opponentShieldHealth -= 10;
//         }
//         else
//         {
//             if (isOpponentVisible)
//             {
//                 opponentActions.BloodSplash();
//             }
//             currentOpponentHealth -= 10;
//             opponentHealthBar.SetOpponentHealth(currentOpponentHealth);
//         }
//         if (currentAmmo == 0)
//         {
//             currentAmmo = maxAmmo;
//             ammoBar.SetAmmo(currentAmmo);
//         }
//     }

//     public void OpponentAmmo()
//     {
//         if (shieldHealth != 0)
//         {
//             shieldHealth -= 10;
//         }
//         else
//         {
//             currentHealth -= 10;
//             healthBar.SetHealth(currentHealth);
//             isDamaged = true;
//         }
//     }

//     public void Grenade()
//     {
//         if (isOpponentVisible)
//         {
//             if (opponentShieldHealth == 10)
//             {
//                 opponentShieldHealth = 0;
//                 currentOpponentHealth -= 20;
//                 opponentHealthBar.SetOpponentHealth(currentOpponentHealth);
//                 opponentActions.BloodSplash();
//             } else if (opponentShieldHealth == 20)
//             {
//                 opponentShieldHealth = 0;
//                 currentOpponentHealth -= 10;
//                 opponentHealthBar.SetOpponentHealth(currentOpponentHealth);
//                 opponentActions.BloodSplash();
//             } else if (opponentShieldHealth == 30)
//             {
//                 opponentShieldHealth = 0;
//             } else
//             {
//                 currentOpponentHealth -= 30;
//                 opponentHealthBar.SetOpponentHealth(currentOpponentHealth);
//                 opponentActions.BloodSplash();
//             }
//         }
        
//         if (currentGrenade == 0)
//         {
//             currentGrenade = maxGrenade;
//             grenadeBar.SetGrenade(currentGrenade);
//         }
//     }

//     public void OpponentGrenade()
//     {
//         if (shieldHealth == 10)
//         {
//             shieldHealth = 0;
//             currentHealth -= 20;
//             healthBar.SetHealth(currentHealth);
//             isDamaged = true;
//         }
//         else if (shieldHealth == 20)
//         {
//             shieldHealth = 0;
//             currentHealth -= 10;
//             healthBar.SetHealth(currentHealth);
//             isDamaged = true;
//         }
//         else if (shieldHealth == 30)
//         {
//             shieldHealth = 0;
//         }
//         else
//         {
//             currentHealth -= 30;
//             healthBar.SetHealth(currentHealth);
//             isDamaged = true;
//         }
//     }

//     public void SetShieldActive()
//     {
//         if (isShieldCooldown == true)
//         {
//            // can't use shield if it's on cooldown
//         } else
//         {
//             currentShield -= 1;
//             shieldUI.SetShield(currentShield);

//             isShieldCooldown = true;
//             shieldCooldownFill = 1;
//             shieldUI.SetShieldCooldownImage(shieldCooldownFill);
//             shieldHealth = 30;
//         }
//     }

//     void Shield()
//     {
//         if (isShieldCooldown)
//         {
//             shieldCooldownFill -= 1 / shieldCooldownTime * Time.deltaTime;
//             shieldUI.SetShieldCooldownImage(shieldCooldownFill);
//             shieldedBorder.SetActive(true);

//             if (shieldCooldownFill <= 0)
//             {
//                 shieldCooldownFill = 0;
//                 shieldUI.SetShieldCooldownImage(shieldCooldownFill);
//                 shieldedBorder.SetActive(false);
//                 isShieldCooldown = false;
//                 shieldHealth = 0;
//             }
//             if (shieldHealth == 0)
//             {
//                 shieldedBorder.SetActive(false);
//             }
//         }
//         if (currentShield == 0)
//         {
//             currentShield = maxShield;
//             shieldUI.SetShield(currentShield);
//         }
//     }

//     public void SetOpponentShieldActive()
//     {
//         if (isOpponentShieldCooldown == true)
//         {
//             // can't use shield if it's on cooldown
//         }
//         else
//         {
//             isOpponentShieldCooldown = true;
//             opponentShieldCooldownFill = 1;
//             //shieldUI.SetShieldCooldownImage(shieldCooldownFill);
//             opponentShieldHealth = 30;
//         }
//     }

//     void OpponentShield()
//     {
//         if (isOpponentShieldCooldown)
//         {
//             opponentShieldCooldownFill -= 1 / opponentShieldCooldownTime * Time.deltaTime;
//             //shieldUI.SetShieldCooldownImage(shieldCooldownFill);
//             opponentShield.SetActive(true);

//             if (opponentShieldCooldownFill <= 0)
//             {
//                 opponentShieldCooldownFill = 0;
//                 //shieldUI.SetShieldCooldownImage(shieldCooldownFill);
//                 opponentShield.SetActive(false);
//                 isOpponentShieldCooldown = false;
//                 opponentShieldHealth = 0;
//             }
//             if (opponentShieldHealth == 0)
//             {
//                 opponentShield.SetActive(false);
//             }
//         }
//         if (currentShield == 0)
//         {
//             currentShield = maxShield;
//             shieldUI.SetShield(currentShield);
//         }
//     }

//     void Damaged()
//     {
//         if (isDamaged)
//         {
//             if (damagedTime > 0)
//             {
//                 damagedBorder.SetActive(true);
//                 damagedTime -= Time.deltaTime;
//             } else
//             {
//                 damagedBorder.SetActive(false);
//                 isDamaged = false;
//                 damagedTime = 1;
//             }
            
//         }
//     }

//     void Score()
//     {
//         player1Score.text = player1.ToString();
//         player2Score.text = player2.ToString();
//     }

//     public void SetOpponentVisible()
//     {
//         isOpponentVisible = true;
//     }

//     public void SetOpponentNotVisible()
//     {
//         isOpponentVisible = false;
//     }
// }
