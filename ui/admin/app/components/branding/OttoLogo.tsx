import { assetUrl, cn } from "~/lib/utils";

import { TypographyH2 } from "~/components/Typography";
import { useTheme } from "~/components/theme";

export function Orion101Logo({
    hideText = false,
    classNames = {},
}: {
    hideText?: boolean;
    classNames?: { root?: string; image?: string };
}) {
    const { isDark } = useTheme();
    let logo = isDark
        ? "/logo/orion101-logo-blue-white-text.svg"
        : "/logo/orion101-logo-blue-black-text.svg";
    if (hideText) {
        logo = "/logo/orion101-icon-blue.svg";
    }
    return (
        <TypographyH2
            className={cn(
                "text-center flex gap-2 items-center justify-center pb-0",
                classNames.root
            )}
        >
            <img
                src={assetUrl(logo)}
                alt="Orion101 Logo"
                className={cn("h-8", classNames.image)}
            />
        </TypographyH2>
    );
}
