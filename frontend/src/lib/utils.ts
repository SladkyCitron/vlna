export { cn } from "cn";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any }
  ? Omit<T, "children">
  : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & {
  ref?: U | null;
};

export function getLocalizedCountryName(countryCode: string): string {
  try {
    const displayNames = new Intl.DisplayNames([navigator.language], {
      type: "region",
    });
    return displayNames.of(countryCode) || countryCode;
  } catch (error) {
    console.error("Error getting localized country name:", error);
    return countryCode;
  }
}

export function getLocalizedLanguageName(languageCode: string): string {
  try {
    const displayNames = new Intl.DisplayNames([navigator.language], {
      type: "language",
    });
    return displayNames.of(languageCode) || languageCode;
  } catch (error) {
    console.error("Error getting localized language name:", error);
    return languageCode;
  }
}
