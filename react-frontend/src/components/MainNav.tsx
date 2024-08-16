"use client";

import * as React from "react";
import Link from "react-router-dom"; 
import { cn } from "@/lib/utils";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import CustomLink from "./custom-link";
import { Button } from "./ui/button";

const categories = [
  {
    title: "Women",
    href: "/women",
    description: "Explore our collection for women."
  },
  {
    title: "Men",
    href: "/men",
    description: "Discover our men's collection."
  },
  {
    title: "Kids",
    href: "/kids",
    description: "Find the latest styles for kids."
  },
  {
    title: "Accessories",
    href: "/accessories",
    description: "Shop our range of accessories."
  },
];

export function MainNav() {
  return (
    <div className="flex gap-6 items-center">
      <CustomLink href="/" className="mr-2">
        <Button variant="ghost" className="p-2 dark:hover:bg-blue-600 rounded-full">
          <img
            src="/vite.svg"
            alt="Home"
            width="32"
            height="32"
            className="min-w-8"
          />
        </Button>
      </CustomLink>
      <NavigationMenu>
        <NavigationMenuList>
          <NavigationMenuItem>
            <NavigationMenuTrigger>Menu</NavigationMenuTrigger>
            <NavigationMenuContent className="absolute top-12 left-0 w-64 bg-white dark:bg-slate-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg">
              <ul className="p-4 space-y-2">
                {categories.map((category) => (
                  <ListItem
                    key={category.title}
                    href={category.href}
                    title={category.title}
                  >
                    {category.description}
                  </ListItem>
                ))}
              </ul>
            </NavigationMenuContent>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>
    </div>
  );
}

const ListItem = React.forwardRef<
  React.ElementRef<"a">,
  React.ComponentPropsWithoutRef<"a">
>(({ className, title, children, ...props }, ref) => {
  return (
    <li>
      <NavigationMenuLink asChild>
        <a
          ref={ref}
          className={cn(
            "block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
            className
          )}
          {...props}
        >
          <div className="text-sm font-medium leading-none">{title}</div>
          <p className="line-clamp-2 text-sm leading-snug text-muted-foreground">
            {children}
          </p>
        </a>
      </NavigationMenuLink>
    </li>
  );
});
ListItem.displayName = "ListItem";
