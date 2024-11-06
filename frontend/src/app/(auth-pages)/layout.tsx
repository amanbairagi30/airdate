import { Badge } from "@/components/ui/badge";
import { coupleGaming, menGaming2, womenGaming2 } from "@/constants/images";
import { Crown } from "lucide-react";
import type { Metadata } from "next";
import Image from "next/image";

export const metadata: Metadata = {
  title: "Airdate - Find Gaming Partners",
  description: "Connect with gamers and content creators",
};

export default function AuthPagesLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <section className="flex items-center w-full h-screen">
      <div className="hidden xl:flex w-[50%] relative bg-gradient-to-t from-primary overflow-hidden to-purple-400 h-screen p-10 flex-col gap-2 justify-end">
        <Image
          className="w-[18rem] opacity-100 rotate-[6deg] absolute top-[-6rem] left-[-5rem] z-1  shadow-2xl shadow-primary/10 rounded-xl h-[26rem]"
          src={womenGaming2}
          alt="men-gmaing"
        />
        <Image
          className="w-[18rem]  absolute rotate-[18deg] bottom-[-13rem] right-[-3rem] z-1   shadow-2xl shadow-primary/50 rounded-xl h-[26rem]"
          src={coupleGaming}
          alt="men-gmaing"
        />
        <Image
          className="w-[18rem]  rotate-[-18deg] absolute top-[-2rem] right-[-8rem] z-1  shadow-2xl shadow-primary/10 rounded-xl h-[26rem]"
          src={menGaming2}
          alt="women-gmaing"
        />
        <div className="flex flex-col gap-4">
          <Badge className="w-fit h-10 z-10">
            <Crown className="w-4 h-4 mr-2" /> #1 platform for{" "}
            <span className="font-serifItalic mx-1 font-normal">match</span>{" "}
            making.
          </Badge>
          <div className="flex flex-col gap-4 z-10">
            <h1 className="text-6xl font-bold">airDate.</h1>
            <p className="w-[70%] text-sm">
              {" "}
              Connect with fellow gamers, share your favorite games, and find
              the perfect gaming partners. Join our community to start
              collaborating and making new friends!
            </p>
          </div>
        </div>
      </div>
      <div className="w-full xl:w-[50%] relative overflow-hidden flex flex-col items-center justify-between md:justify-center py-10 h-screen">
        <div className="w-full mx-2 space-y-10 px-4">{children}</div>
        <div className="text-center m-4 p-4 xl:hidden">
          <h2 className="text-3xl font-bold tracking-tight font-serifItalic">
            <span className="">air</span>
            <span className="text-purple-400">Date .</span>
          </h2>
          <p className="mt-2 max-w-xs mx-auto text-sm text-muted-foreground">
            Connect with gamers and content creators and find your perfect match
          </p>
        </div>
        <div className="absolute z-[-1] transition-all duration-700 ease-out opacity-45 bottom-[-30rem] left-[50%] translate-x-[-50%] rounded-xl bg-gradient-to-t from-purple-400 to-purple-700 blur-[6em] size-[40rem]"></div>
      </div>
    </section>
  );
}
