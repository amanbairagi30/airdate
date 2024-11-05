import React from "react";
import { Button } from "./ui/button";
import { Badge } from "./ui/badge";
import Image from "next/image";
import {
  coupleGaming,
  menGaming,
  menGaming2,
  womenGaming,
  womenGaming2,
} from "@/constants/images";
import { Heart, Telescope } from "lucide-react";

export default function Hero() {
  return (
    <>
      <section className="w-full relative rotate-12min-h-[80vh] max-h-fit flex flex-col items-center justify-start my-20">
        <div className="absolute z-[-1] transition-all duration-700 ease-out opacity-35 -left-[15rem] rounded-xl bg-gradient-to-t from-purple-400 to-purple-700 blur-[8em] -top-[35rem] size-[40rem]"></div>
        <div className="absolute transition-all duration-700 ease-out opacity-35 -right-[15rem] rounded-xl bg-gradient-to-t from-purple-400 to-purple-700 blur-[8em] top-[25rem] size-[40rem]"></div>

        <Badge className="p-1 rounded-full text-xs px-4 bg-accent/50 text-foreground mb-4 hover:bg-accent/50 border-2 border-primary/80">
          Being Loved by many folks ✨
        </Badge>
        <div className="container px-4 md:px-6">
          <div className="flex flex-col items-center justify-center space-y-8 text-center">
            <div className="space-y-2 flex flex-col gap-4">
              <h1 className="text-3xl font-semibold max-w-4xl sm:text-4xl md:text-5xl lg:text-6xl/none">
                Meet your
                <span className="bg-gradient-to-tr pl-3 font-serifItalic from-primary to-purple-300 bg-clip-text text-transparent">
                  perfect
                </span>{" "}
                partner in
                <div className="mt-2 font-semibold">gaming and creation.</div>
              </h1>
              <p className="mx-auto max-w-[600px] font-medium md:text-lg">
                Connect with fellow gamers and content creators who'll make your
                heart rate and APM skyrocket.
              </p>
            </div>
            <div className="flex items-center justify-center gap-4">
              <Button className="">
                Find my match <Telescope className="w-4 h-4" />
              </Button>
              <Button variant={"linkHover2"} className="">
                Try for free
              </Button>
            </div>
          </div>
        </div>

        <section className="my-20 w-full flex flex-col gap-6 items-center justify-center">
          <div className="flex relative items-center gap-24">
            <Image
              className="w-[18rem] -rotate-[6deg]  shadow-2xl shadow-primary/10 rounded-xl h-[26rem]"
              src={menGaming2}
              alt="men-gmaing"
            />
            <Image
              className="w-[18rem] absolute top-[50%] left-[50%] translate-x-[-50%] z-20 shadow-2xl shadow-primary/50 translate-y-[-50%] rounded-xl h-[26rem]"
              src={coupleGaming}
              alt="men-gmaing"
            />
            <Image
              className="w-[18rem] rotate-[6deg]  shadow-2xl shadow-primary/10 rounded-xl h-[26rem]"
              src={womenGaming2}
              alt="women-gmaing"
            />
          </div>

          <div className="mt-10 border-[3px] bg-gradient-to-tl from-primary/30 to-purple-300/10 border-primary flex items-center px-2 gap-2 w-[14rem] rounded-full h-[4rem]">
            <div className="bg-gradient-to-t from-primary w-10 h-10 rounded-full to-purple-300 flex items-center justify-center">
              <Heart className="text-purple-300 fill-purple-300" />
            </div>
            <div className="font-serifItalic flex-1 text-center">
              Love is in the air
            </div>
          </div>
        </section>
      </section>
    </>
  );
}
