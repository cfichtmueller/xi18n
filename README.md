# xi18n

Generate `Localizable.strings` for Xcode projects.

## Usage

1. Initialize xi18n in your project directory: `xi18n init`
2. Modify `x18n.yml` to suit your needs
3. Generate assets by running `xi18n gen`

## Project File Format

```yaml
languages:
    de: [Project]/de.lproj/Localizable.strings
    en: [Project]/en.lproj/Localizable.strings
stringsFile: [Project]/Strings.swift
keys:
    Hello:
        de: Hallo
        en: Hello
    World:
        de: Welt
        en: World

```

- Replace `[Project]` with the name of your project.
- Create one key per language. `xi18n` will create `Localizable.strings` for each defined language.
- Specify the path of the generated `Strings.swift` file.
- Specify the keys. Keys support an optional `comment`.

## Usage in Code

Pretty Simple.

```swift
import SwiftUI

struct MyView: View {
    var body: some View {
        HStack {
            Text(Strings.Hello)
            Text(Strings.World)
        }
    }
}
```