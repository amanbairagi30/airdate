import React from "react";

export default function Footer() {
  return (
    <>
      <footer className="h-fit py-8 bg-gradient-to-b from-accent/50  rounded-t-3xl">
        <section className="max-w-7xl mx-auto">
          <div className="flex items-center justify-between">
            <div className="text-3xl font-secondary font-bold">
              Pixel<span className="text-purple-400">&</span>Chill
            </div>
            All rights reserved @ {new Date().getFullYear()}
          </div>
        </section>
      </footer>
    </>
  );
}
