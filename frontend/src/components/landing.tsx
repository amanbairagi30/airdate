import React from "react";
import { NavBar } from "./NavBar";
import { SelectTheme } from "./theme-toggle";
import Hero from "./hero";
import Features from "./features";
import Footer from "./footer";

export default function Landing() {
  return (
    <>
      <main className="max-w-7xl p-4 mx-auto">
        <section className="p-4">
          <NavBar />
          <Hero />
          <Features />
        </section>
      </main>
      <Footer />
    </>
  );
}
