"use client";

import { cn } from "@/lib/utils";
import CustomLink from "./custom-link";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import React from "react";
import { Button } from "./ui/button";
import { Link } from 'react-router-dom';

export function MainNav() {
  return (
    <div className="flex gap-6 items-center p-2">
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
            <NavigationMenuLink
              href="/product"
              className={cn(navigationMenuTriggerStyle(), "bg-blue-500 dark:text-gray-200 dark:bg-primary hover:bg-gray-200 dark:hover:bg-blue-600 px-4 py-2 rounded-md")}
            >
              Product
            </NavigationMenuLink>
          </NavigationMenuItem>
          <NavigationMenuItem>
            <NavigationMenuLink
              href="/about"
              className={cn(navigationMenuTriggerStyle(), "bg-blue-500 dark:text-gray-200 dark:bg-primary hover:bg-gray-200 dark:hover:bg-blue-600 px-4 py-2 rounded-md")}
            >
              About
            </NavigationMenuLink>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>
    </div>
  );
}

const ListItem = React.forwardRef<
  HTMLAnchorElement,
  { to: string; title: string; className?: string; children?: React.ReactNode }
>(({ to, title, children, ...props }, ref) => {
  return (
    <li>
      <NavigationMenuLink asChild>
        <Link
          to={to}
          ref={ref}
          className={cn(
            "block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
            props.className
          )}
          {...props}
        >
          <div className="text-sm font-medium leading-none">{title}</div>
          <p className="text-sm leading-snug line-clamp-2 text-muted-foreground">
            {children}
          </p>
        </Link>
      </NavigationMenuLink>
    </li>
  );
});
ListItem.displayName = "ListItem";
