import { QueryFunctionContext } from "@tanstack/react-query"

export interface ConfigurationPage {
  // ported from the [go struct](https://github.com/rehomed/homepage/blob/main/config.go#L21)
  title?: string;
  path: string;
  inject: {
    css?: string;
    js?: string;
    html?: string;
  };
  search?: ConfigurationPageSearch;
  widgets?: ConfigurationPageWidget[];
  links?: ConfigurationPageLink[];
}

export interface ConfigurationPageSearch {
  enabled?: boolean;
  label?: string;
  url?: string;
}

export interface ConfigurationPageWidget {
  builtin: string;
  inject: string;
  kv: Record<string, string>;
}

export interface ConfigurationPageLink {
  title?: string;
  color?: string;
  icon?: string;
  url?: string;
}

export async function getConfig({ queryKey }: QueryFunctionContext) {
  const path = queryKey[0]
  const req = await fetch(import.meta.env.VITE_PUBLIC_API_BASE + path + ".json")

  if (req.status !== 200) {
    throw new Error(`failed with status ${req.status}: ${await req.text()}`)
  }

  const json = await req.json()
  return json as ConfigurationPage
}