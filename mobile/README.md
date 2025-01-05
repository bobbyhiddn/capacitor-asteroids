# Mobile Deployment Guide

## Prerequisites

### iOS
1. Mac computer with latest macOS
2. Xcode (latest version)
3. Apple Developer Account ($99/year)
4. iOS Developer Certificate
5. App Store Connect account

### Android
1. Android Studio
2. Google Play Developer Account ($25 one-time fee)
3. Java Development Kit (JDK)
4. Android SDK

## Building the Apps

### iOS Build Steps
1. Install Xcode from the Mac App Store
2. Sign in to your Apple Developer account in Xcode
3. Run the following commands:
```bash
cd mobile/ios
gomobile bind -target=ios -o Game.xcframework github.com/bobbyhiddn/ecs-asteroids
```
4. Open the Xcode project and configure:
   - Bundle identifier
   - Version number
   - Build number
   - App icons
   - Launch screen
   - Required device capabilities
   - Privacy descriptions

### Android Build Steps
1. Install Android Studio
2. Set ANDROID_HOME environment variable
3. Run the following commands:
```bash
cd mobile/android
gomobile bind -target=android -o game.aar github.com/bobbyhiddn/ecs-asteroids
```
4. Configure in Android Studio:
   - Package name
   - Version code and name
   - App icons
   - Permissions
   - Minimum SDK version

## App Store Submission Checklist

### iOS App Store
1. Screenshots (6.5" & 5.5" displays required)
2. App description
3. Keywords
4. Privacy policy URL
5. Support URL
6. Marketing URL (optional)
7. App Store icon (1024x1024)
8. App Store rating questionnaire
9. Export compliance information
10. Content rights declaration

### Google Play Store
1. Screenshots (phone & tablet)
2. Feature graphic (1024x500)
3. App description
4. Privacy policy URL
5. Content rating questionnaire
6. Pricing & distribution settings
7. App category
8. Contact information
9. Store listing localization (optional)

## Important Notes
- Test thoroughly on multiple devices before submission
- Follow platform-specific design guidelines
- Implement appropriate privacy measures
- Keep certificates and signing keys secure
- Plan for app updates and maintenance
