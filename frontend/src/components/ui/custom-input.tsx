import { Eye, EyeOff, Info, Key, Mail } from "lucide-react";
import { useState } from "react";
import { Field, FieldLabel } from "./field";
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
  description,
}: {
  displayName: string;
  labelProps: React.ComponentProps<"label">;
  inputProps: React.ComponentProps<"input">;
  textAddon?: string;
  description?: string;
}) {
  return (
    <Field>
      <FieldLabel className="flex flex-col items-start gap-0" {...labelProps}>
        <div className="">
          {displayName}
          {inputProps?.required && <span className="ms-1 text-red">*</span>}
        </div>
      </FieldLabel>
      <InputGroup>
        <InputGroupInput {...inputProps} autoComplete="off" />
        <InputGroupAddon align="inline-end">
          {textAddon && <InputGroupText>{textAddon}</InputGroupText>}
        </InputGroupAddon>
      </InputGroup>
      {description && (
        <div className="flex items-center gap-0.5 ps-2 text-amber">
          <Info className="size-3.5" />
          <div className="text-xs font-medium">{description}</div>
        </div>
      )}
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
  defaultValue,
}: {
  displayName: string;
  items: { label: string; value: string }[];
  required?: boolean;
  onValueChange: (value: unknown) => void;
  defaultValue?: unknown;
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
        defaultValue={defaultValue}
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

export function PasswordInput({
  label,
  id,
  name,
  className,
  description,
  ...props
}: { label: string; description?: string } & React.ComponentProps<"input">) {
  const [show, setShow] = useState<boolean>(false);
  return (
    <div className="grid gap-2">
      <label htmlFor={id} className="text-sm flex items-center gap-3">
        {label}
      </label>
      <InputGroup>
        <InputGroupAddon align="inline-start">
          <Key />
        </InputGroupAddon>
        <InputGroupInput
          id={id}
          name={name}
          type={show ? "text" : "password"}
          autoComplete="current-password"
          placeholder="********"
          required
          {...props}
        />
        <InputGroupAddon
          align="inline-end"
          onClick={() => setShow((prevState) => !prevState)}
          className="hover:cursor-pointer hover:text-blue"
        >
          {show ? <EyeOff /> : <Eye />}
        </InputGroupAddon>
      </InputGroup>
      {description && (
        <div className="flex items-center gap-0.5 ps-2 text-amber">
          <Info className="size-3.5" />
          <div className="text-xs font-medium">{description}</div>
        </div>
      )}
    </div>
  );
}

export function EmailInput({
  label,
  id,
  name,
  className,
  ...props
}: { label: string } & React.ComponentProps<"input">) {
  return (
    <div className="grid gap-3">
      <label htmlFor={id} className="text-sm">
        {label}
      </label>
      <div>
        <InputGroup>
          <InputGroupAddon align="inline-start">
            <Mail />
          </InputGroupAddon>
          <InputGroupInput
            placeholder="mail@kmed.com"
            id={id}
            name={name}
            autoComplete="off"
            type="email"
            {...props}
          />
        </InputGroup>
      </div>
    </div>
  );
}
