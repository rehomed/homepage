import { VNode } from "preact";
import {
  ConfigurationPageLink,
  ConfigurationPageSearch,
  ConfigurationPageWidget,
  getConfig,
} from "../../config";
import { useEffect, useState } from "preact/hooks";
import { useQuery } from "@tanstack/react-query";
import clsx from "clsx";
import { hexToRGBA } from "../../util";

function Layout({ children }: { children: VNode[] | VNode }) {
  return (
    <div className="w-screen h-screen flex flex-col items-center px-2 md:p-0">
      <div className="max-w-2xl w-full grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
        {children}
      </div>
    </div>
  );
}

// order:
// - search
// - links

function Search({
  config,
  isSkeleton,
}: {
  config?: ConfigurationPageSearch;
  isSkeleton?: boolean;
}) {
  if (isSkeleton) {
    return <div className="w-full h-12 skeleton shadow-xl" />;
  }

  if (!config || !config.enabled) return null;

  const onKeyDown = (e: KeyboardEvent & { target: HTMLInputElement }) => {
    if (e.key === "Enter") {
      const searchURL = String(config.url).replace(/%s/, e?.target?.value);

      if (e.ctrlKey) {
        window.open(searchURL);
      } else {
        window.location.href = searchURL;
      }
    }
  };

  return (
    <input
      placeholder={
        config.label ? `Search with ${config.label}` : "Search query"
      }
      className="input input-md bg-zinc-950/25 w-full shadow-xl"
      autoFocus
      onKeyDown={onKeyDown as (e: KeyboardEvent) => void}
      id="search"
    ></input>
  );
}

function Link({
  config,
  isSkeleton,
}: {
  config?: ConfigurationPageLink;
  isSkeleton?: boolean;
}) {
  if (isSkeleton) {
    return (
      <div
        className={clsx(
          "w-full h-20 skeleton rounded-none",
          "border-l-2 border-l-zinc-700 shadow"
        )}
      />
    );
  }

  if (!config) return null;

  const hostname = config.url ? new URL(config.url).host : "";
  const color = config.color || "#ffffff";
  const colorRGBA = hexToRGBA(color, 0.3);

  return (
    <div className="bg-zinc-900/35 shadow">
      <a
        href={config.url}
        target="_blank"
        rel="noopener noreferrer"
        className={clsx(
          "border-l-2 bg-gradient-to-r",
          "p-4 w-full h-full",
          "flex flex-col items-start",
          "hover:scale-[105%] duration-300"
        )}
        style={{
          borderLeft: "2px solid " + color,
          background: `linear-gradient(to right, ${colorRGBA}, transparent)`,
        }}
      >
        <div className="flex flex-row gap-2 items-center">
          <img
            // TODO get favicons on the server
            src={
              config.icon ||
              `https://external-content.duckduckgo.com/ip3/${hostname}.ico`
            }
            className="w-6 h-6 shrink-0"
          />
          <h1 className="font-semibold text-lg">{config.title}</h1>
        </div>
        <h1 className="mt-1 text-xs text-zinc-400">{hostname}</h1>
      </a>
    </div>
  );
}

function Widget({ config }: { config?: ConfigurationPageWidget }) {
  if (!config) return null;

  return (
    <>
      <iframe
        frameborder={0}
        src={config.inject}
        className={clsx("bg-transparent w-full")}
      ></iframe>
    </>
  );
}

export default function Main() {
  const [path, setPath] = useState("");
  useEffect(() => {
    if (typeof window !== "undefined") {
      setPath(window.location.pathname);
    }
  }, [window]);

  const { isLoading, error, data } = useQuery({
    queryKey: [path],
    queryFn: getConfig,
  });

  if (error) {
    return (
      <Layout>
        <div
          className={clsx(
            "p-4 md:col-span-2  overflow-x-scroll",
            "bg-red-950/10 border-l-2 border-l-red-200/50"
          )}
        >
          <h1 className="text-lg font-semibold">Something went wrong:</h1>
          <div className="collapse bg-zinc-950/25">
            <input type="checkbox" />
            <div className="collapse-title">
              <h2 className="text-lg">{error.message}</h2>
              <h4 className="text-xs text-zinc-300">expand for more details</h4>
            </div>
            <div className="collapse-content">
              <pre className="text-red-400">{error.stack}</pre>
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  if (isLoading) {
    return (
      <Layout>
        <div className="md:col-span-2 w-full">
          <Search isSkeleton />
        </div>
        {/* this is stupid. idc. */}
        <Link isSkeleton />
        <Link isSkeleton />
        <Link isSkeleton />
        <Link isSkeleton />
        <Link isSkeleton />
        <Link isSkeleton />
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="md:col-span-2 w-full">
        <Search config={data?.search} />
      </div>
      <>
        {/* typescipt screams out in pain about this for some reason */}
        {data?.widgets
          ? data.widgets.map((w: ConfigurationPageWidget) => (
              <Widget config={w} />
            ))
          : null}
        {data?.links
          ? data.links.map((l: ConfigurationPageLink) => <Link config={l} />)
          : null}
      </>
    </Layout>
  );
}
