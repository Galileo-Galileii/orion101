import { cn } from "~/lib/utils";

import { TypographyP } from "~/components/Typography";
import {
    Tooltip,
    TooltipContent,
    TooltipTrigger,
} from "~/components/ui/tooltip";

export function TruncatedText({
    content,
    className,
}: {
    content: React.ReactNode;
    className?: string;
}) {
    return (
        <Tooltip>
            <TooltipTrigger asChild>
                <div className={cn(`truncate cursor-pointer`, className)}>
                    <TypographyP className="truncate">{content}</TypographyP>
                </div>
            </TooltipTrigger>
            <TooltipContent>
                <p>{content}</p>
            </TooltipContent>
        </Tooltip>
    );
}
