// ocr_darwin.m - Objective-C implementation for macOS Vision OCR

#import <Foundation/Foundation.h>
#import <Vision/Vision.h>
#import <CoreImage/CoreImage.h>
#include <stdlib.h>

// Recognize text from image at path using Vision framework.
// Returns a C string containing the recognized text, joined by newlines.
// The caller is responsible for freeing the returned string.
char* recognizeText(const char* imagePath) {
    @autoreleasepool {
        NSString *path = [NSString stringWithUTF8String:imagePath];
        NSURL *imageURL = [NSURL fileURLWithPath:path];

        // Load image
        CIImage *image = [CIImage imageWithContentsOfURL:imageURL];
        if (!image) {
            return NULL;
        }

        VNImageRequestHandler *handler = [[VNImageRequestHandler alloc] initWithCIImage:image options:@{}];

        // Create request
        VNRecognizeTextRequest *request = [[VNRecognizeTextRequest alloc] initWithCompletionHandler:nil];
        request.recognitionLevel = VNRequestTextRecognitionLevelAccurate;
        request.usesLanguageCorrection = YES;
        // automaticallyDetectsLanguage requires macOS 13.0+
        if (@available(macOS 13.0, *)) {
            request.automaticallyDetectsLanguage = YES;
        }
        // Prioritize common languages: Chinese (Simplified/Traditional), English, Japanese, Korean, etc.
        request.recognitionLanguages = @[@"zh-Hans", @"zh-Hant", @"en-US", @"ja-JP", @"ko-KR", @"de-DE", @"fr-FR", @"es-ES"];

        NSError *error = nil;
        [handler performRequests:@[request] error:&error];

        if (error) {
            return NULL;
        }

        NSMutableString *resultText = [NSMutableString string];
        for (VNRecognizedTextObservation *observation in request.results) {
            VNRecognizedText *text = [observation topCandidates:1].firstObject;
            if (text) {
                if (resultText.length > 0) {
                    [resultText appendString:@"\n"];
                }
                [resultText appendString:text.string];
            }
        }

        return strdup([resultText UTF8String]);
    }
}
