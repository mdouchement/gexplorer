// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin && !ios
// +build darwin,!ios

//
// https://developer.apple.com/documentation/appkit/nsopenpanel
// https://stackoverflow.com/questions/1640419/open-file-dialog-box (objective-c)
// https://github.com/tcltk/tk/blob/b44ff01fb12579f2d87626ab15f41bcef0d0609b/macosx/tkMacOSXDialog.c#L377-L407
// https://github.com/electron/electron/blob/main/shell/browser/ui/file_dialog_mac.mm
// https://ourcodeworld.com/articles/read/1117/how-to-implement-a-file-and-directory-picker-in-macos-using-swift-5
//

#include "_cgo_export.h"
#import <Foundation/Foundation.h>
#import <Appkit/AppKit.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

void exportFile(CFTypeRef viewRef, int32_t id, char * name) {
	NSView *view = (__bridge NSView *)viewRef;

	NSSavePanel *panel = [NSSavePanel savePanel];

    [panel setNameFieldStringValue:@(name)];
	[panel beginSheetModalForWindow:[view window] completionHandler:^(NSModalResponse result){ // FIXME: NSSavePanel: 0x100890620> running implicitly; please run panels using NSSavePanel rather than NSApplication.
		if (result == NSModalResponseOK) {
			exportCallback(id, (char *)[[[panel URL] absoluteString] UTF8String]);
		} else {
		    exportCallback(id, (char *)(""));
		}
	}];
}

void importFile(CFTypeRef viewRef, int32_t id, char * ext) {
    NSMutableArray<NSString*> *exts = [[@(ext) componentsSeparatedByString:@","] mutableCopy];
    NSMutableArray<UTType*> *contentTypes = [[NSMutableArray alloc]init];

    int i;
    for (i = 0; i < [exts count]; i++) {
        UTType * utt = [UTType typeWithFilenameExtension:exts[i]];
        if (utt != nil){
            [contentTypes addObject:utt];
        }
     }

	NSOpenPanel *panel = [NSOpenPanel openPanel];
    [panel setAllowedContentTypes:[NSArray arrayWithArray:contentTypes]];

	NSView *view = (__bridge NSView *)viewRef;
	[panel beginSheetModalForWindow:[view window] completionHandler:^(NSModalResponse result){ // FIXME: NSOpenPanel: 0x100989b00> running implicitly; please run panels using NSSavePanel rather than NSApplication.
		if (result == NSModalResponseOK) {
			importCallback(id, (char *)[[[panel URL] absoluteString] UTF8String]);
		} else {
			importCallback(id, (char *)(""));
		}
	}];
}

void importFiles(CFTypeRef viewRef, int32_t id, char * ext) {
	NSView *view = (__bridge NSView *)viewRef;

	NSOpenPanel *panel = [NSOpenPanel openPanel];
	// [panel setCanChooseFiles:YES];
	// [panel setCanChooseDirectories:NO];
	[panel setAllowsMultipleSelection:YES];

    NSMutableArray<NSString*> *exts = [[@(ext) componentsSeparatedByString:@","] mutableCopy];
    NSMutableArray<UTType*> *contentTypes = [[NSMutableArray alloc]init];

    int i;
    for (i = 0; i < [exts count]; i++) {
        UTType * utt = [UTType typeWithFilenameExtension:exts[i]];
        if (utt != nil){
            [contentTypes addObject:utt];
        }
     }
    [panel setAllowedContentTypes:[NSArray arrayWithArray:contentTypes]];

	[panel beginSheetModalForWindow:[view window] completionHandler:^(NSModalResponse result){ // FIXME: NSOpenPanel: 0x100989b00> running implicitly; please run panels using NSSavePanel rather than NSApplication.
		if (result == NSModalResponseOK) {
			NSArray* urls = [panel URLs];
			NSInteger count = [urls count];

			char* results[count];
			for(int i = 0; i < count; i++)	{
				results[i] = (char *)[[[urls objectAtIndex:i] absoluteString] UTF8String];
			}

			// for(int i = 0; i < count; i++)
			// {
			// 	NSLog(@"URL: %s", results[i]);
			// }

			importsCallback(id, count, results);
		} else {
		    importCallback(id, (char *)("")); // Use the single import to ease the implementation.
		}
	}];
}