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

/**
 * English (US) translations for Thunder applications
 * All namespaces organized in a single file for better maintainability
 */
const translations = {
  // ============================================================================
  // Common namespace - Shared translations across all Thunder applications
  // ============================================================================
  common: {
    // Actions
    actions: {
      add: 'Add',
      edit: 'Edit',
      delete: 'Delete',
      cancel: 'Cancel',
      save: 'Save',
      create: 'Create',
      update: 'Update',
      remove: 'Remove',
      search: 'Search',
      filter: 'Filter',
      reset: 'Reset',
      submit: 'Submit',
      close: 'Close',
      back: 'Back',
      next: 'Next',
      previous: 'Previous',
      confirm: 'Confirm',
      ok: 'OK',
      yes: 'Yes',
      no: 'No',
      continue: 'Continue',
      skip: 'Skip',
      finish: 'Finish',
      refresh: 'Refresh',
      copy: 'Copy',
      download: 'Download',
      upload: 'Upload',
      export: 'Export',
      import: 'Import',
      view: 'View',
      details: 'Details',
      settings: 'Settings',
      logout: 'Logout',
      login: 'Login',
    },

    // Status messages
    status: {
      loading: 'Loading...',
      saving: 'Saving...',
      success: 'Success',
      error: 'Error',
      warning: 'Warning',
      info: 'Info',
      pending: 'Pending',
      active: 'Active',
      inactive: 'Inactive',
      enabled: 'Enabled',
      disabled: 'Disabled',
      completed: 'Completed',
      failed: 'Failed',
    },

    // Form labels
    form: {
      name: 'Name',
      description: 'Description',
      email: 'Email',
      password: 'Password',
      username: 'Username',
      required: 'Required',
      optional: 'Optional',
      requiredField: 'This field is required',
      invalidEmail: 'Invalid email address',
      invalidFormat: 'Invalid format',
      searchPlaceholder: 'Search...',
    },

    // Messages
    messages: {
      confirmDelete: 'Are you sure you want to delete this item?',
      deleteSuccess: 'Item deleted successfully',
      deleteError: 'Failed to delete item',
      saveSuccess: 'Saved successfully',
      saveError: 'Failed to save',
      updateSuccess: 'Updated successfully',
      updateError: 'Failed to update',
      createSuccess: 'Created successfully',
      createError: 'Failed to create',
      noData: 'No data available',
      noResults: 'No results found',
      somethingWentWrong: 'Something went wrong',
      tryAgain: 'Please try again',
    },

    // Navigation
    navigation: {
      home: 'Home',
      dashboard: 'Dashboard',
      profile: 'Profile',
      help: 'Help',
      documentation: 'Documentation',
    },

    // User menu
    userMenu: {
      profile: 'Profile',
      myAccount: 'My account',
      addAnotherAccount: 'Add another account',
      settings: 'Settings',
      signOut: 'Sign Out',
    },

    // Header
    header: {
      notifications: 'Coming soon',
      openNotifications: 'Open notifications',
    },

    // Data table - MUI DataGrid locale text
    dataTable: {
      // Root
      noRowsLabel: 'No rows',
      noResultsOverlayLabel: 'No results found.',
      noColumnsOverlayLabel: 'No columns',
      noColumnsOverlayManageColumns: 'Manage columns',

      // Density selector toolbar button text
      toolbarDensity: 'Density',
      toolbarDensityLabel: 'Density',
      toolbarDensityCompact: 'Compact',
      toolbarDensityStandard: 'Standard',
      toolbarDensityComfortable: 'Comfortable',

      // Columns selector toolbar button text
      toolbarColumns: 'Columns',
      toolbarColumnsLabel: 'Select columns',

      // Filters toolbar button text
      toolbarFilters: 'Filters',
      toolbarFiltersLabel: 'Show filters',
      toolbarFiltersTooltipHide: 'Hide filters',
      toolbarFiltersTooltipShow: 'Show filters',
      toolbarFiltersTooltipActive: (count: number) =>
        count !== 1 ? `${count} active filters` : `${count} active filter`,

      // Quick filter toolbar field
      toolbarQuickFilterPlaceholder: 'Search…',
      toolbarQuickFilterLabel: 'Search',
      toolbarQuickFilterDeleteIconLabel: 'Clear',

      // Export selector toolbar button text
      toolbarExport: 'Export',
      toolbarExportLabel: 'Export',
      toolbarExportCSV: 'Download as CSV',
      toolbarExportPrint: 'Print',

      // Columns management text
      columnsManagementSearchTitle: 'Search',
      columnsManagementNoColumns: 'No columns',
      columnsManagementShowHideAllText: 'Show/Hide All',
      columnsManagementReset: 'Reset',

      // Filter panel text
      filterPanelAddFilter: 'Add filter',
      filterPanelRemoveAll: 'Remove all',
      filterPanelDeleteIconLabel: 'Delete',
      filterPanelLogicOperator: 'Logic operator',
      filterPanelOperator: 'Operator',
      filterPanelOperatorAnd: 'And',
      filterPanelOperatorOr: 'Or',
      filterPanelColumns: 'Columns',
      filterPanelInputLabel: 'Value',
      filterPanelInputPlaceholder: 'Filter value',

      // Filter operators text
      filterOperatorContains: 'contains',
      filterOperatorDoesNotContain: 'does not contain',
      filterOperatorEquals: 'equals',
      filterOperatorDoesNotEqual: 'does not equal',
      filterOperatorStartsWith: 'starts with',
      filterOperatorEndsWith: 'ends with',
      filterOperatorIs: 'is',
      filterOperatorNot: 'is not',
      filterOperatorAfter: 'is after',
      filterOperatorOnOrAfter: 'is on or after',
      filterOperatorBefore: 'is before',
      filterOperatorOnOrBefore: 'is on or before',
      filterOperatorIsEmpty: 'is empty',
      filterOperatorIsNotEmpty: 'is not empty',
      filterOperatorIsAnyOf: 'is any of',

      // Filter values text
      filterValueAny: 'any',
      filterValueTrue: 'true',
      filterValueFalse: 'false',

      // Column menu text
      columnMenuLabel: 'Menu',
      columnMenuShowColumns: 'Show columns',
      columnMenuManageColumns: 'Manage columns',
      columnMenuFilter: 'Filter',
      columnMenuHideColumn: 'Hide column',
      columnMenuUnsort: 'Unsort',
      columnMenuSortAsc: 'Sort by ASC',
      columnMenuSortDesc: 'Sort by DESC',

      // Column header text
      columnHeaderFiltersTooltipActive: (count: number) =>
        count !== 1 ? `${count} active filters` : `${count} active filter`,
      columnHeaderFiltersLabel: 'Show filters',
      columnHeaderSortIconLabel: 'Sort',

      // Rows selected footer text
      footerRowSelected: (count: number) =>
        count !== 1 ? `${count.toLocaleString()} rows selected` : `${count.toLocaleString()} row selected`,

      // Total row amount footer text
      footerTotalRows: 'Total Rows:',

      // Total visible row amount footer text
      footerTotalVisibleRows: (visibleCount: number, totalCount: number) =>
        `${visibleCount.toLocaleString()} of ${totalCount.toLocaleString()}`,

      // Checkbox selection text
      checkboxSelectionHeaderName: 'Checkbox selection',
      checkboxSelectionSelectAllRows: 'Select all rows',
      checkboxSelectionUnselectAllRows: 'Unselect all rows',
      checkboxSelectionSelectRow: 'Select row',
      checkboxSelectionUnselectRow: 'Unselect row',

      // Boolean cell text
      booleanCellTrueLabel: 'yes',
      booleanCellFalseLabel: 'no',

      // Actions cell more text
      actionsCellMore: 'more',

      // Column pinning text
      pinToLeft: 'Pin to left',
      pinToRight: 'Pin to right',
      unpin: 'Unpin',

      // Tree Data
      treeDataGroupingHeaderName: 'Group',
      treeDataExpand: 'see children',
      treeDataCollapse: 'hide children',

      // Grouping columns
      groupingColumnHeaderName: 'Group',
      groupColumn: (name: string) => `Group by ${name}`,
      unGroupColumn: (name: string) => `Stop grouping by ${name}`,

      // Master/detail
      detailPanelToggle: 'Detail panel toggle',
      expandDetailPanel: 'Expand',
      collapseDetailPanel: 'Collapse',

      // Pagination
      paginationRowsPerPage: 'Rows per page:',
      paginationDisplayedRows: ({from, to, count}: {from: number; to: number; count: number}) =>
        `${from}–${to} of ${count !== -1 ? count : `more than ${to}`}`,

      // Row reordering text
      rowReorderingHeaderName: 'Row reordering',

      // Aggregation
      aggregationMenuItemHeader: 'Aggregation',
      aggregationFunctionLabelSum: 'sum',
      aggregationFunctionLabelAvg: 'avg',
      aggregationFunctionLabelMin: 'min',
      aggregationFunctionLabelMax: 'max',
      aggregationFunctionLabelSize: 'size',
    },
  },

  // ============================================================================
  // Navigation namespace - Navigation related translations
  // ============================================================================
  navigation: {
    pages: {
      users: 'Users',
      userTypes: 'User Types',
      integrations: 'Integrations',
      applications: 'Applications',
      dashboard: 'Dashboard',
    },
    breadcrumb: {
      develop: 'Develop',
    },
  },

  // ============================================================================
  // Users namespace - User management feature translations
  // ============================================================================
  users: {
    title: 'User Management',
    subtitle: 'Manage users, roles, and permissions across your organization',
    addUser: 'Add User',
    editUser: 'Edit User',
    deleteUser: 'Delete User',
    userDetails: 'User Details',
    firstName: 'First Name',
    lastName: 'Last Name',
    email: 'Email Address',
    username: 'Username',
    role: 'Role',
    status: 'Status',
    createdAt: 'Created At',
    lastLogin: 'Last Login',
    actions: 'Actions',
    noUsers: 'No users found',
    searchUsers: 'Search users...',
    confirmDeleteUser: 'Are you sure you want to delete this user?',
    userCreatedSuccess: 'User created successfully',
    userUpdatedSuccess: 'User updated successfully',
    userDeletedSuccess: 'User deleted successfully',
  },

  // ============================================================================
  // User Types namespace - User types feature translations
  // ============================================================================
  userTypes: {
    title: 'User Types',
    addUserType: 'Add User Type',
    createUserType: 'Create User Type',
    editUserType: 'Edit User Type',
    deleteUserType: 'Delete User Type',
    userTypeDetails: 'User Type Details',
    typeName: 'Type Name',
    typeNamePlaceholder: 'e.g., Employee, Customer, Partner',
    description: 'Description',
    createDescription: 'Define a new user type schema for your organization',
    permissions: 'Permissions',
    schemaProperties: 'Schema Properties',
    propertyName: 'Property Name',
    propertyNamePlaceholder: 'e.g., email, age, address',
    propertyType: 'Type',
    addProperty: 'Add Property',
    unique: 'Unique',
    regexPattern: 'Regular Expression Pattern (Optional)',
    regexPlaceholder: 'e.g., ^[a-zA-Z0-9]+$',
    enumValues: 'Allowed Values (Enum) - Optional',
    enumPlaceholder: 'Add value and press Enter',
    types: {
      string: 'String',
      number: 'Number',
      boolean: 'Boolean',
      enum: 'Enum',
      array: 'Array',
      object: 'Object',
    },
    validationErrors: {
      nameRequired: 'Please enter a user type name',
      propertiesRequired: 'Please add at least one property',
      duplicateProperties: 'Duplicate property names found: {{duplicates}}',
    },
    noUserTypes: 'No user types found',
    confirmDeleteUserType: 'Are you sure you want to delete this user type?',
  },

  // ============================================================================
  // Integrations namespace - Integrations feature translations
  // ============================================================================
  integrations: {
    title: 'Integrations',
    subtitle: 'Manage your integrations and connections',
    addIntegration: 'Add Integration',
    editIntegration: 'Edit Integration',
    deleteIntegration: 'Delete Integration',
    integrationDetails: 'Integration Details',
    provider: 'Provider',
    apiKey: 'API Key',
    endpoint: 'Endpoint',
    status: 'Status',
    connected: 'Connected',
    disconnected: 'Disconnected',
    testConnection: 'Test Connection',
    noIntegrations: 'No integrations found',
    comingSoon: 'Coming Soon',
    comingSoonDescription: 'Integrations management functionality will be available soon.',
  },

  // ============================================================================
  // Applications namespace - Applications feature translations
  // ============================================================================
  applications: {
    title: 'Applications',
    subtitle: 'Manage your applications and services',
    addApplication: 'Add Application',
    editApplication: 'Edit Application',
    deleteApplication: 'Delete Application',
    applicationDetails: 'Application Details',
    appName: 'Application Name',
    appType: 'Application Type',
    clientId: 'Client ID',
    clientSecret: 'Client Secret',
    redirectUri: 'Redirect URI',
    allowedOrigins: 'Allowed Origins',
    noApplications: 'No applications found',
    comingSoon: 'Coming Soon',
    comingSoonDescription: 'Applications management functionality will be available soon.',
  },

  // ============================================================================
  // Dashboard namespace - Dashboard feature translations
  // ============================================================================
  dashboard: {
    welcomeMessage: 'Welcome to Thunder Develop',
    totalUsers: 'Total Users',
    activeUsers: 'Active Users',
    totalApplications: 'Total Applications',
    recentActivity: 'Recent Activity',
    quickActions: 'Quick Actions',
  },

  // ============================================================================
  // Authentication namespace - Authentication feature translations
  // ============================================================================
  auth: {
    signIn: 'Sign In',
    signUp: 'Sign Up',
    signOut: 'Sign Out',
    forgotPassword: 'Forgot Password?',
    resetPassword: 'Reset Password',
    changePassword: 'Change Password',
    rememberMe: 'Remember Me',
    welcomeBack: 'Welcome Back',
    createAccount: 'Create Account',
    alreadyHaveAccount: 'Already have an account?',
    dontHaveAccount: "Don't have an account?",
    enterEmail: 'Enter your email',
    enterPassword: 'Enter your password',
    confirmPassword: 'Confirm Password',
    passwordMismatch: 'Passwords do not match',
    invalidCredentials: 'Invalid credentials',
    accountLocked: 'Account is locked',
    sessionExpired: 'Session expired. Please sign in again.',
    signInSuccess: 'Signed in successfully',
    signUpSuccess: 'Account created successfully',
    passwordResetSent: 'Password reset link sent to your email',
    passwordResetSuccess: 'Password reset successfully',
  },

  // ============================================================================
  // MFA namespace - Multi-factor authentication feature translations
  // ============================================================================
  mfa: {
    title: 'Multi-Factor Authentication',
    setupMfa: 'Set Up MFA',
    enableMfa: 'Enable MFA',
    disableMfa: 'Disable MFA',
    verificationCode: 'Verification Code',
    enterCode: 'Enter verification code',
    sendCode: 'Send Code',
    resendCode: 'Resend Code',
    invalidCode: 'Invalid verification code',
    codeExpired: 'Verification code expired',
    scanQrCode: 'Scan this QR code with your authenticator app',
    backupCodes: 'Backup Codes',
    saveBackupCodes: 'Save these backup codes in a secure place',
  },

  // ============================================================================
  // Social namespace - Social login feature translations
  // ============================================================================
  social: {
    continueWith: 'Continue with',
    signInWith: 'Sign in with',
    google: 'Google',
    facebook: 'Facebook',
    github: 'GitHub',
    microsoft: 'Microsoft',
  },

  // ============================================================================
  // Consent namespace - Consent feature translations
  // ============================================================================
  consent: {
    title: 'Consent Required',
    message: 'The application is requesting access to your information',
    requestedPermissions: 'Requested Permissions',
    allow: 'Allow',
    deny: 'Deny',
    learnMore: 'Learn More',
  },

  // ============================================================================
  // Errors namespace - Error messages translations
  // ============================================================================
  errors: {
    authenticationFailed: 'Authentication failed',
    unauthorizedAccess: 'Unauthorized access',
    accessDenied: 'Access denied',
    invalidRequest: 'Invalid request',
    serverError: 'Server error occurred',
    networkError: 'Network error. Please check your connection.',
    redirectFailed: 'Redirect failed',
  },
} as const;

export default translations;
