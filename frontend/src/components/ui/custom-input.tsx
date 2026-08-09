import { Field, FieldLabel } from "./field";
import { Input } from "./input";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupTextarea,
} from "./input-group";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./select";

/*
 * Custom Input for use in this form
 */
export function CInput({
  displayName,
  labelProps,
  inputProps,
  textAddon,
}: {
  displayName: string;
  labelProps: React.ComponentProps<"label">;
  inputProps: React.ComponentProps<"input">;
  textAddon?: string;
}) {
  return (
    <Field>
      <FieldLabel {...labelProps}>
        {displayName}{" "}
        {inputProps?.required && <span className="text-red">*</span>}
      </FieldLabel>
      <InputGroup>
        <InputGroupInput {...inputProps} autoComplete="off" />
        <InputGroupAddon align="inline-end">
          {textAddon && <InputGroupText>{textAddon}</InputGroupText>}
        </InputGroupAddon>
      </InputGroup>
    </Field>
  );
}
export function CTextArea({
  displayName,
  labelProps,
  textareaProps,
}: {
  displayName: string;
  labelProps: React.ComponentProps<"label">;
  textareaProps: React.ComponentProps<"textarea">;
}) {
  return (
    <Field>
      <FieldLabel {...labelProps}>
        {displayName}{" "}
        {textareaProps?.required && <span className="text-red">*</span>}
      </FieldLabel>
      <InputGroup>
        <InputGroupTextarea {...textareaProps} />
      </InputGroup>
    </Field>
  );
}
export function CSelectInput({
  displayName,
  items,
  required,
  onValueChange,
}: {
  displayName: string;
  items: { label: string; value: string }[];
  required?: boolean;
  onValueChange: (value: unknown) => void;
}) {
  return (
    <Field>
      <FieldLabel>
        {displayName}
        {required && <span className="text-red">*</span>}
      </FieldLabel>
      <Select
        items={items}
        autoComplete="off"
        onValueChange={onValueChange}
        required
      >
        <SelectTrigger>
          <SelectValue />
        </SelectTrigger>
        <SelectContent alignItemWithTrigger={false}>
          <SelectGroup>
            {items.map((item) => (
              <SelectItem key={item.label} value={item.value}>
                {item.label}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </Field>
  );
}
