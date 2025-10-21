/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import {Controller} from 'react-hook-form';
import type {Control, FieldErrors, Path} from 'react-hook-form';
import {useState} from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import FormLabel from '@mui/material/FormLabel';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Chip from '@mui/material/Chip';
import IconButton from '@mui/material/IconButton';
import {Plus} from 'lucide-react';
import type {PropertyDefinition} from '../types/users';

/**
 * Array input component for adding multiple values as chips
 */
function ArrayFieldInput({
  value,
  onChange,
  fieldLabel,
}: {
  value: string[];
  onChange: (value: string[]) => void;
  fieldLabel: string;
}) {
  const [inputValue, setInputValue] = useState('');
  const currentValue = Array.isArray(value) ? value : [];

  const handleAdd = () => {
    if (inputValue.trim()) {
      onChange([...currentValue, inputValue.trim()]);
      setInputValue('');
    }
  };

  const handleDelete = (indexToDelete: number) => {
    onChange(currentValue.filter((_, index) => index !== indexToDelete));
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAdd();
    }
  };

  return (
    <Box>
      <Box sx={{display: 'flex', gap: 1, mb: 1}}>
        <TextField
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={`Add ${fieldLabel.toLowerCase()}`}
          fullWidth
          size="small"
          variant="outlined"
        />
        <IconButton size="small" onClick={handleAdd} disabled={!inputValue.trim()}>
          <Plus size={16} />
        </IconButton>
      </Box>
      <Box sx={{display: 'flex', flexWrap: 'wrap', gap: 1}}>
        {currentValue.length > 0 &&
          currentValue.map((item) => {
            const itemIndex = currentValue.indexOf(item);
            return (
              <Chip
                key={`chip-${item}-${Math.random()}`}
                label={String(item)}
                onDelete={() => handleDelete(itemIndex)}
                variant="outlined"
                size="medium"
              />
            );
          })}
      </Box>
    </Box>
  );
}

/**
 * Helper function to convert field name to a readable label
 * Examples: firstname -> Firstname, partnershipType -> Partnership Type
 */
export const formatFieldLabel = (fieldName: string): string => {
  // Handle camelCase by adding spaces before capital letters
  const spacedName = fieldName.replace(/([A-Z])/g, ' $1').trim();
  // Capitalize first letter
  return spacedName.charAt(0).toUpperCase() + spacedName.slice(1);
};

/**
 * Helper function to render a form field based on the property definition
 *
 * @param fieldName - The name of the field in the schema
 * @param fieldDef - The property definition from the schema
 * @param control - React Hook Form control object
 * @param errors - Form validation errors
 * @returns A rendered form field component or null for unsupported types
 */
export const renderSchemaField = <T extends Record<string, unknown>>(
  fieldName: string,
  fieldDef: PropertyDefinition,
  control: Control<T>,
  errors: FieldErrors<T>,
) => {
  const isRequired = fieldDef.required ?? false;
  const fieldLabel = formatFieldLabel(fieldName);

  // String fields
  if (fieldDef.type === 'string') {
    const stringDef = fieldDef;

    // Render as Select dropdown if enum values are provided
    if (stringDef.enum && stringDef.enum.length > 0) {
      const enumOptions = stringDef.enum;
      return (
        <FormControl key={fieldName}>
          <FormLabel htmlFor={fieldName}>
            {fieldLabel}
            {isRequired && <span style={{color: 'red'}}> *</span>}
          </FormLabel>
          <Controller
            name={fieldName as Path<T>}
            control={control}
            rules={{
              required: isRequired ? `${fieldLabel} is required` : false,
            }}
            render={({field}) => (
              <Select
                {...field}
                value={field.value ?? ''}
                id={fieldName}
                fullWidth
                required={isRequired}
                error={!!errors[fieldName]}
                displayEmpty
              >
                <MenuItem value="">
                  <em>Select {fieldLabel}</em>
                </MenuItem>
                {enumOptions.map((option) => (
                  <MenuItem key={option} value={option}>
                    {option.charAt(0).toUpperCase() + option.slice(1)}
                  </MenuItem>
                ))}
              </Select>
            )}
          />
          {errors[fieldName] && (
            <Typography variant="caption" color="error" sx={{mt: 0.5, ml: 1.75}}>
              {errors[fieldName]?.message as string}
            </Typography>
          )}
        </FormControl>
      );
    }

    // Render as TextField for regular string fields
    // Determine validation pattern
    let validationPattern;
    if (stringDef.regex) {
      validationPattern = {
        value: new RegExp(stringDef.regex),
        message: `${fieldLabel} format is invalid`,
      };
    }

    return (
      <FormControl key={fieldName}>
        <FormLabel htmlFor={fieldName}>
          {fieldLabel}
          {isRequired && <span style={{color: 'red'}}> *</span>}
        </FormLabel>
        <Controller
          name={fieldName as Path<T>}
          control={control}
          rules={{
            required: isRequired ? `${fieldLabel} is required` : false,
            pattern: validationPattern,
          }}
          render={({field}) => (
            <TextField
              {...field}
              id={fieldName}
              type={stringDef.type}
              placeholder={`Enter ${fieldLabel.toLowerCase()}`}
              fullWidth
              required={isRequired}
              variant="outlined"
              error={!!errors[fieldName]}
              helperText={errors[fieldName]?.message as string}
              color={errors[fieldName] ? 'error' : 'primary'}
            />
          )}
        />
      </FormControl>
    );
  }

  // Number fields
  if (fieldDef.type === 'number') {
    return (
      <FormControl key={fieldName}>
        <FormLabel htmlFor={fieldName}>
          {fieldLabel}
          {isRequired && <span style={{color: 'red'}}> *</span>}
        </FormLabel>
        <Controller
          name={fieldName as Path<T>}
          control={control}
          rules={{
            required: isRequired ? `${fieldLabel} is required` : false,
          }}
          render={({field}) => (
            <TextField
              {...field}
              id={fieldName}
              type="number"
              placeholder={`Enter ${fieldLabel.toLowerCase()}`}
              fullWidth
              required={isRequired}
              variant="outlined"
              error={!!errors[fieldName]}
              helperText={errors[fieldName]?.message as string}
              color={errors[fieldName] ? 'error' : 'primary'}
              onChange={(e) => {
                const {value} = e.target;
                field.onChange(value ? Number(value) : '');
              }}
            />
          )}
        />
      </FormControl>
    );
  }

  // Boolean fields
  if (fieldDef.type === 'boolean') {
    return (
      <FormControl key={fieldName}>
        <Controller
          name={fieldName as Path<T>}
          control={control}
          render={({field}) => (
            <Box sx={{display: 'flex', alignItems: 'center', py: 1}}>
              <input
                type="checkbox"
                id={fieldName}
                name={field.name}
                checked={!!field.value}
                onChange={(e) => field.onChange(e.target.checked)}
                onBlur={field.onBlur}
                ref={field.ref}
                style={{width: '18px', height: '18px', cursor: 'pointer'}}
              />
              <FormLabel htmlFor={fieldName} sx={{ml: 1.5, cursor: 'pointer', mb: 0}}>
                {fieldLabel}
                {isRequired && <span style={{color: 'red'}}> *</span>}
              </FormLabel>
            </Box>
          )}
        />
      </FormControl>
    );
  }

  // Array fields
  if (fieldDef.type === 'array') {
    return (
      <FormControl key={fieldName} fullWidth>
        <FormLabel htmlFor={fieldName}>
          {fieldLabel}
          {isRequired && <span style={{color: 'red'}}> *</span>}
        </FormLabel>
        <Controller
          name={fieldName as Path<T>}
          control={control}
          rules={{
            required: isRequired ? `${fieldLabel} is required` : false,
            validate: (value) => {
              if (isRequired && (!Array.isArray(value) || value.length === 0)) {
                return `${fieldLabel} must have at least one value`;
              }
              return true;
            },
          }}
          render={({field}) => {
            const fieldValue = field.value as string[];
            return (
              <Box>
                <ArrayFieldInput value={fieldValue} onChange={field.onChange} fieldLabel={fieldLabel} />
                {errors[fieldName] && (
                  <Typography variant="caption" color="error" sx={{mt: 0.5, ml: 1.75}}>
                    {errors[fieldName]?.message as string}
                  </Typography>
                )}
              </Box>
            );
          }}
        />
      </FormControl>
    );
  }

  // For unsupported types, return null
  return null;
};
