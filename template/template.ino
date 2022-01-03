// Program im Debugmodus kompilieren, dann werden zus. Ausgaben auf die serielle Schnittstelle geschrieben.
//#define debug
/*
   Here are the defines used in this software to control special parts of the implementation
   #define TPS_RCRECEIVER: using a RC receiver input
   #define TPS_ENHANCEMENT: all of the other enhancments
   #define TPS_SERVO: using servo outputs
   #define TPS_TONE: using a tone output

*/
// defining different hardware platforms
#ifdef __AVR_ATmega328P__
//#define TPS_RCRECEIVER
#define TPS_ENHANCEMENT
#define TPS_SERVO
#define TPS_TONE
#endif

#ifdef ESP32
//#define TPS_RCRECEIVER (not implementted yet)
//#define TPS_ENHANCEMENT
//#define TPS_SERVO
//#define TPS_TONE
#endif

#ifdef __AVR_ATtiny84__
//#define TPS_ENHANCEMENT
//#define TPS_SERVO
//#define TPS_TONE
#endif

#ifdef __AVR_ATtiny861__
//#define TPS_RCRECEIVER
//#define TPS_ENHANCEMENT
//#define TPS_SERVO
//#define TPS_TONE
#endif

#ifdef __AVR_ATtiny4313__
// because of the limited memory only 2 of this three options are available. 
//#define TPS_RCRECEIVER
//#define TPS_ENHANCEMENT
//#define TPS_TONE
#endif

// libraries
#include "debug.h"
#include "makros.h"

#ifdef ESP32
#include <ESP32Servo.h>
#ifdef TPS_TONE
#include <ESP32Tone.h>
#endif
#endif

#ifdef TPS_SERVO
#if defined(__AVR_ATmega328P__) || defined(__AVR_ATtiny84__) || defined(__AVR_ATtiny861__) || defined(__AVR_ATtiny4313__)
#include <Servo.h>
#endif
#endif

#ifdef TPS_TONE
#include "notes.h"
#endif

#include "hardware.h"
// sub routines
const byte subCnt = 7;
word subs[subCnt];

// the actual address of the program
word addr;
// page register
word page;
// defining register
byte a, b, c, d;
#ifdef TPS_ENHANCEMENT
byte e, f;
#endif

#ifdef TPS_ENHANCEMENT
const byte SAVE_CNT = 16;
#else
const byte SAVE_CNT = 1;
#endif

word saveaddr[SAVE_CNT];
int saveCnt;

#ifdef TPS_ENHANCEMENT
byte stack[SAVE_CNT];
byte stackCnt;
#endif

unsigned long tmpValue;

#ifdef TPS_SERVO
Servo servo1;
Servo servo2;
#endif

void setup() {
  // put your setup code here, to run once:
  pinMode(Dout_1, OUTPUT);
  pinMode(Dout_2, OUTPUT);
  pinMode(Dout_3, OUTPUT);
  pinMode(Dout_4, OUTPUT);

  pinMode(PWM_1, OUTPUT);
  pinMode(PWM_2, OUTPUT);

  pinMode(Din_1, INPUT_PULLUP);
  pinMode(Din_2, INPUT_PULLUP);
  pinMode(Din_3, INPUT_PULLUP);
  pinMode(Din_4, INPUT_PULLUP);

  pinMode(SW_PRG, INPUT_PULLUP);
  pinMode(SW_SEL, INPUT_PULLUP);

  initHardware();

  digitalWrite(Dout_1, 1);
  delay(1000);
  digitalWrite(Dout_1, 0);

  // Serielle Schnittstelle einstellen
  initDebug();

  doReset();

#ifdef TPS_ENHANCEMENT
  pinMode(LED_BUILTIN, OUTPUT);
#endif

// todo setup servos
{{.setup}}
}

void loop() {
  // put your main code here, to run repeatedly:
{{.main}}
}

void doReset() {
  dbgOutLn("Reset");
#ifdef TPS_SERVO
  servo1.detach();
  servo2.detach();
#endif

  for (int i = 0; i < subCnt; i++) {
    subs[i] = 0;
  }

  addr = 0;
  page = 0;
  saveCnt = 0;
  a = 0;
  b = 0;
  c = 0;
  d = 0;
#ifdef TPS_ENHANCEMENT
  e = 0;
  f = 0;
  stackCnt = 0;
  for (int i = 0; i < SAVE_CNT; i++) {
    stack[i] = 0;
  }
#endif
}

/*
  output to port
*/
void doPort(byte data) {
  digitalWrite(Dout_1, (data & 0x01) > 0);
  digitalWrite(Dout_2, (data & 0x02) > 0);
  digitalWrite(Dout_3, (data & 0x04) > 0);
  digitalWrite(Dout_4, (data & 0x08) > 0);
}
