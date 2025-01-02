import { Airplay, ArrowRightIcon, Columns2, Film, Users } from "lucide-react";
import React from "react";
import { Button } from "./ui/button";

export default function Features() {
  return (
    <section className="mt-32">
      <div className="flex items-center justify-center">
        <div className="text-center">
          <div className="text-xl font-semibold">Here are some</div>
          <div className="text-4xl mt-1 font-serifItalic">Latest Offerings</div>
        </div>
      </div>

      <div className="w-full grid grid-cols-1 md:grid-cols-2 my-16 gap-4">
        <div className="h-[25rem] bg-gradient-to-t transition-all duration-700 ease-in-out from-accent/50 hover:shadow-2xl group p-4 overflow-hidden relative rounded-xl border">
          <div className="absolute transition-all duration-700 ease-out opacity-0 group-hover:opacity-55 -left-[15rem] rounded-full bg-gradient-to-t from-purple-400 to-purple-700 blur-[8em] -top-[35rem] size-[40rem]"></div>

          <div className="flex absolute top-[-0rem] opacity-25 group-hover:opacity-80 right-[-2rem] transition-all duration-500 ease-in-out group-hover:rotate-[-16deg] -rotate-12 items-center justify-center p-2 border-2 rounded-xl w-fit bg-primary">
            <div className="w-40 h-40 flex items-center justify-center rounded-lg bg-gradient-to-b from-purple-400 to-purple-600">
              <Users className="w-16 h-16 text-white" />
            </div>
          </div>

          <div className="flex flex-col gap-2 absolute bottom-6">
            <h3 className="text-3xl font-semibold">
              Find <span className="font-serifItalic">Matches</span>
            </h3>
            <p className="text-sm text-gray-500 font-semibold">
              Swipe through profiles to find players and creators who share your
              passion for gaming and content creation.
            </p>
            <Button
              Icon={ArrowRightIcon}
              iconPlacement="right"
              className="w-fit mt-4 text-xs h-6"
            >
              Read more
            </Button>
          </div>
        </div>
        <div className="h-[25rem] bg-gradient-to-t transition-all duration-700 ease-in-out from-accent/50 hover:shadow-2xl group p-4 overflow-hidden relative rounded-xl border">
          <div className="absolute transition-all duration-700 ease-out opacity-0 group-hover:opacity-55 -left-[15rem] rounded-full bg-gradient-to-t from-purple-400 to-purple-700 blur-[8em] -top-[35rem] size-[40rem]"></div>

          <div className="flex absolute top-[-0rem] opacity-25 group-hover:opacity-80 right-[-2rem] transition-all duration-500 ease-in-out group-hover:rotate-[-16deg] -rotate-12 items-center justify-center p-2 border-2 rounded-xl w-fit bg-primary">
            <div className="w-40 h-40 flex items-center justify-center rounded-lg bg-gradient-to-b from-purple-400 to-purple-600">
              <Columns2 className="w-16 h-16 text-white" />
            </div>
          </div>

          <div className="flex flex-col gap-2 absolute bottom-6">
            <h3 className="text-3xl font-semibold">
              Duo <span className="font-serifItalic">Queue</span>
            </h3>
            <p className="text-sm text-gray-500 font-semibold">
              Looking for a gaming buddy? Get paired with someone ready to team
              up for epic sessions or video collabs.
            </p>
            <Button
              Icon={ArrowRightIcon}
              iconPlacement="right"
              className="w-fit mt-4 text-xs h-6"
            >
              Read more
            </Button>
          </div>
        </div>
        <div className="h-[25rem] bg-gradient-to-t transition-all duration-700 ease-in-out from-accent/50 hover:shadow-2xl group p-4 overflow-hidden relative rounded-xl border">
          <div className="absolute transition-all duration-700 ease-out opacity-0 group-hover:opacity-55 -left-[15rem] rounded-full bg-gradient-to-t from-purple-400 to-purple-700 blur-[8em] -top-[35rem] size-[40rem]"></div>

          <div className="flex absolute top-[-0rem] opacity-25 group-hover:opacity-80 right-[-2rem] transition-all duration-500 ease-in-out group-hover:rotate-[-16deg] -rotate-12 items-center justify-center p-2 border-2 rounded-xl w-fit bg-primary">
            <div className="w-40 h-40 flex items-center justify-center rounded-lg bg-gradient-to-b from-purple-400 to-purple-600">
              <Airplay className="w-16 h-16 text-white" />
            </div>
          </div>

          <div className="flex flex-col gap-2 absolute bottom-6">
            <h3 className="text-3xl font-semibold">
              Raid my <span className="font-serifItalic">Stream</span>
            </h3>
            <p className="text-sm text-gray-500 font-semibold">
              Invite matched players to join your live stream as special guests.
              More fun, more hype, and more viewers!
            </p>
            <Button
              Icon={ArrowRightIcon}
              iconPlacement="right"
              className="w-fit mt-4 text-xs h-6"
            >
              Read more
            </Button>
          </div>
        </div>
        <div className="h-[25rem] bg-gradient-to-t transition-all duration-700 ease-in-out from-accent/50 hover:shadow-2xl group p-4 overflow-hidden relative rounded-xl border">
          <div className="absolute transition-all duration-700 ease-out opacity-0 group-hover:opacity-55 -left-[15rem] rounded-full bg-gradient-to-t from-purple-400 to-purple-700 blur-[8em] -top-[35rem] size-[40rem]"></div>

          <div className="flex absolute top-[-0rem] opacity-25 group-hover:opacity-80 right-[-2rem] transition-all duration-500 ease-in-out group-hover:rotate-[-16deg] -rotate-12 items-center justify-center p-2 border-2 rounded-xl w-fit bg-primary">
            <div className="w-40 h-40 flex items-center justify-center rounded-lg bg-gradient-to-b from-purple-400 to-purple-600">
              <Film className="w-16 h-16 text-white" />
            </div>
          </div>

          <div className="flex flex-col gap-2 absolute bottom-6">
            <h3 className="text-3xl font-semibold">
              Content <span className="font-serifItalic">Collab</span>
            </h3>
            <p className="text-sm text-gray-500 font-semibold">
              Get ideas for content collabs based on your style and interests.
              Make memorable videos together.
            </p>
            <Button
              Icon={ArrowRightIcon}
              iconPlacement="right"
              className="w-fit mt-4 text-xs h-6"
            >
              Read more
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}
