// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Foundation/Foundation.h>
#import <objc/runtime.h>
#include <AppKit/AppKit.h>
#include <errno.h>

@implementation NSBundle(noti)
- (NSString *)notiIdentifier {
	return @"com.apple.terminal";
}
@end

@interface NotiDelegate : NSObject<NSUserNotificationCenterDelegate>
@property (nonatomic, assign) BOOL delivered;
@end

@implementation NotiDelegate
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didActivateNotification:(NSUserNotification *)notification {
	self.delivered = YES;
}
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didDeliverNotification:(NSUserNotification *)notification {
	[NSApp activateIgnoringOtherApps:YES];
	self.delivered = YES;
}
@end

void DesktopNotify(const char* title, const char* message, const char* sound) {
	errno = 0;
	@autoreleasepool {
		Class nsBundle = objc_getClass("NSBundle");
		method_exchangeImplementations(
			class_getInstanceMethod(nsBundle, @selector(bundleIdentifier)),
			class_getInstanceMethod(nsBundle, @selector(notiIdentifier))
		);

		NotiDelegate *notiDel = [[NotiDelegate alloc]init];
		notiDel.delivered = NO;

		NSUserNotificationCenter *nc = [NSUserNotificationCenter defaultUserNotificationCenter];
		nc.delegate = notiDel;

		NSUserNotification *nt = [[NSUserNotification alloc] init];
		nt.title = [NSString stringWithUTF8String:title];
		nt.informativeText = [NSString stringWithUTF8String:message];
		nt.soundName = NSUserNotificationDefaultSoundName;

		if ([[NSString stringWithUTF8String:sound] length] != 0) {
			nt.soundName = [NSString stringWithUTF8String:sound];
		}

		[nc deliverNotification:nt];

		while (notiDel.delivered == NO) {
			[[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
		}
	}
}
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"unsafe"
)

// desktopNotify triggers an NSUserNotification.
func desktopNotify() {
	runUtility()

	sound := os.Getenv(soundEnv)
	if sound == "" {
		sound = "Ping"
	}

	t := C.CString(*title)
	m := C.CString(*message)
	s := C.CString(sound)
	defer C.free(unsafe.Pointer(t))
	defer C.free(unsafe.Pointer(m))
	defer C.free(unsafe.Pointer(s))

	C.DesktopNotify(t, m, s)
}

// speechNotify triggers an NSSpeechSynthesizer notification.
func speechNotify() {
	runUtility()

	voice := os.Getenv(voiceEnv)
	if voice == "" {
		voice = "Alex"
	}
	*message = fmt.Sprintf("%s %s", *title, *message)

	cmd := exec.Command("say", "-v", voice, *message)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
