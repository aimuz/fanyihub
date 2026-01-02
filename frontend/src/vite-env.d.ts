/// <reference types="svelte" />
/// <reference types="vite/client" />

// Wails runtime types
interface WailsRuntime {
  EventsOn(eventName: string, callback: (...data: unknown[]) => void): () => void
  EventsOff(eventName: string): void
  EventsOnce(eventName: string, callback: (...data: unknown[]) => void): () => void
  EventsOnMultiple(
    eventName: string,
    callback: (...data: unknown[]) => void,
    maxCallbacks: number
  ): () => void
  EventsEmit(eventName: string, ...data: unknown[]): void
  WindowReload(): void
  WindowReloadApp(): void
  WindowSetAlwaysOnTop(b: boolean): void
  WindowSetSystemDefaultTheme(): void
  WindowSetLightTheme(): void
  WindowSetDarkTheme(): void
  WindowCenter(): void
  WindowSetTitle(title: string): void
  WindowFullscreen(): void
  WindowUnfullscreen(): void
  WindowIsFullscreen(): Promise<boolean>
  WindowSetSize(width: number, height: number): void
  WindowGetSize(): Promise<{ w: number; h: number }>
  WindowSetMaxSize(width: number, height: number): void
  WindowSetMinSize(width: number, height: number): void
  WindowSetPosition(x: number, y: number): void
  WindowGetPosition(): Promise<{ x: number; y: number }>
  WindowHide(): void
  WindowShow(): void
  WindowMaximise(): void
  WindowToggleMaximise(): void
  WindowUnmaximise(): void
  WindowIsMaximised(): Promise<boolean>
  WindowMinimise(): void
  WindowUnminimise(): void
  WindowIsMinimised(): Promise<boolean>
  WindowIsNormal(): Promise<boolean>
  WindowSetBackgroundColour(R: number, G: number, B: number, A: number): void
  ScreenGetAll(): Promise<Screen[]>
  WindowPrint(): void
  BrowserOpenURL(url: string): void
  Environment(): Promise<EnvironmentInfo>
  Quit(): void
  Hide(): void
  Show(): void
  ClipboardGetText(): Promise<string>
  ClipboardSetText(text: string): Promise<boolean>
  OnFileDrop(
    callback: (x: number, y: number, paths: string[]) => void,
    useDropTarget: boolean
  ): void
  OnFileDropOff(): void
  CanResolveFilePaths(): boolean
  ResolveFilePaths(files: File[]): Promise<string[]>
}

interface Screen {
  isCurrent: boolean
  isPrimary: boolean
  width: number
  height: number
}

interface EnvironmentInfo {
  buildType: string
  platform: string
  arch: string
}

declare global {
  interface Window {
    runtime: WailsRuntime
    go: {
      main: {
        App: {
          GetProviders(): Promise<import('./types').Provider[]>
          AddProvider(provider: import('./types').Provider): Promise<void>
          UpdateProvider(oldName: string, provider: import('./types').Provider): Promise<void>
          RemoveProvider(name: string): Promise<void>
          SetProviderActive(name: string): Promise<void>
          GetActiveProvider(): Promise<import('./types').Provider | null>
          TranslateWithLLM(request: import('./types').TranslateRequest): Promise<string>
          DetectLanguage(text: string): Promise<import('./types').DetectLanguageResponse>
          GetDefaultLanguages(): Promise<Record<string, string>>
          SetDefaultLanguage(sourceLang: string, targetLang: string): Promise<void>
          ToggleWindowVisibility(): Promise<void>
        }
      }
    }
  }
}

export {}
