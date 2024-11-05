import type { Metadata } from "next";
import localFont from "next/font/local";
import {
  Instrument_Sans,
  Instrument_Serif,
  Pixelify_Sans,
  Playfair_Display,
} from "next/font/google";
import "./globals.css";
import { NavBar } from "../components/NavBar";
import { ErrorBoundary } from "../components/ErrorBoundary";
import { ThemeProvider } from "@/components/theme-provider";
import { cn } from "@/lib/utils";

const serifItalic = localFont({
  variable: "--font-serif",
  display: "swap",
  src: [
    {
      path: "./fonts/serif-italic.woff2",
      weight: "400",
      style: "medium",
    },
  ],
});

const instSans = Instrument_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-primary",
  display: "swap",
});

const instSerif = Playfair_Display({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-secondary",
  display: "swap",
});

const pixel = Pixelify_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-pixel",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Airdate - Find Gaming Partners",
  description: "Connect with gamers and content creators",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body
        className={cn(
          instSans.variable,
          pixel.variable,
          instSerif.variable,
          serifItalic.variable,
          "font-primary overflow-x-hidden"
        )}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <ErrorBoundary>{children}</ErrorBoundary>
        </ThemeProvider>
      </body>
    </html>
  );
}
