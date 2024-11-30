import { EyeOff } from "lucide-react";

import { ConfirmationDialog } from "~/components/composed/ConfirmationDialog";
import { Button } from "~/components/ui/button";

type UnpublishProps = {
    className?: string;
    onUnpublish: () => void;
};

export function Unpublish({ onUnpublish }: UnpublishProps) {
    return (
        <ConfirmationDialog
            title="Unpublish Agent"
            description="Are you sure you want to unpublish this agent? This action will disrupt every user currently using this reference."
            onConfirm={() => onUnpublish()}
            confirmProps={{
                variant: "destructive",
                children: "Unpublish",
            }}
        >
            <Button variant="secondary" size="sm">
                <EyeOff className="w-4 h-4" />
                Unpublish
            </Button>
        </ConfirmationDialog>
    );
}
