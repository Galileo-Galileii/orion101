import { toast } from "sonner";
import { mutate } from "swr";

import { OAuthProvider } from "~/lib/model/oauthApps/oauth-helpers";
import { OauthAppService } from "~/lib/service/api/oauthAppService";

import { ConfirmationDialog } from "~/components/composed/ConfirmationDialog";
import { LoadingSpinner } from "~/components/ui/LoadingSpinner";
import { Button } from "~/components/ui/button";
import {
    Tooltip,
    TooltipContent,
    TooltipTrigger,
} from "~/components/ui/tooltip";
import { useOAuthAppInfo } from "~/hooks/oauthApps/useOAuthApps";
import { useAsync } from "~/hooks/useAsync";

export function DeleteOAuthApp({
    id,
    disableTooltip,
    type,
}: {
    id: string;
    disableTooltip?: boolean;
    type: OAuthProvider;
}) {
    const spec = useOAuthAppInfo(type);

    const deleteOAuthApp = useAsync(async () => {
        await OauthAppService.deleteOauthApp(id);
        await mutate(OauthAppService.getOauthApps.key());

        toast.success(`${spec.displayName} OAuth configuration deleted`);
    });

    return (
        <div className="flex gap-2">
            <Tooltip open={getIsOpen()}>
                <ConfirmationDialog
                    title={`Reset ${spec.displayName} OAuth to use Acorn Gateway`}
                    description={`By clicking \`Reset\`, you will delete your custom ${spec.displayName} OAuth configuration and reset to use Acorn Gateway.`}
                    onConfirm={deleteOAuthApp.execute}
                    confirmProps={{
                        variant: "destructive",
                        children: "Reset",
                    }}
                >
                    <TooltipTrigger asChild>
                        <Button
                            variant="destructive"
                            className="w-full"
                            disabled={deleteOAuthApp.isLoading}
                        >
                            {deleteOAuthApp.isLoading ? (
                                <LoadingSpinner className="w-4 h-4 mr-2" />
                            ) : null}
                            Reset {spec.displayName} OAuth to use Acorn Gateway
                        </Button>
                    </TooltipTrigger>
                </ConfirmationDialog>

                <TooltipContent>Delete</TooltipContent>
            </Tooltip>
        </div>
    );

    function getIsOpen() {
        if (disableTooltip) return false;
    }
}