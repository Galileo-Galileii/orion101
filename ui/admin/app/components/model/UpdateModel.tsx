import { Model } from "~/lib/model/models";

import { ModelForm } from "~/components/model/ModelForm";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogTitle,
    DialogTrigger,
} from "~/components/ui/dialog";

type UpdateModelDialogProps = {
    model: Nullish<Model>;
    open: boolean;
    setOpen: (open: boolean) => void;
    children?: React.ReactNode;
};

export function UpdateModelDialog(props: UpdateModelDialogProps) {
    const { model, open, setOpen, children } = props;

    if (!model) return null;

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogTitle>Update model</DialogTitle>

                <DialogDescription hidden>Update model</DialogDescription>

                <ModelForm model={model} onSubmit={() => setOpen(false)} />
            </DialogContent>

            {children && <DialogTrigger asChild>{children}</DialogTrigger>}
        </Dialog>
    );
}
