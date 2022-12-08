--
-- ~/.config/nvim/init.lua
--
----

--
-- I used two references to build up this "solution":
--
--     https://vonheikemen.github.io/devlog/tools/manage-neovim-lsp-client-without-plugins/
--
-- Then for the completion I added the "config.on_attach" section, via this comment:
--
--     https://www.reddit.com/r/neovim/comments/rs47cx/tsserver_and_vimlspomnifunc/
--
-- I'm sure there are better approaches.


--
-- When *.yal files are loaded then set the filetype to be lisp
--
vim.filetype.add {
   pattern = {
      ['.*.yal'] = 'lisp',
   }
}

--
-- Define a helper function which will associate our LSP
-- magic with the appropriate filenames.
--
-- It'll setup completion too.
--
local launch_yal_server = function()
   local autocmd


   local filetypes = {
      'lisp',
   }

   local config = {
      cmd = {'yal', '-lsp'},
      name = 'yal',
      root_dir = vim.fn.getcwd(),
      capabilities = vim.lsp.protocol.make_client_capabilities(),
   }

   -- This gives completion.
   config.on_attach = function(client, bufnr)
      vim.api.nvim_buf_set_option(bufnr, 'omnifunc', 'v:lua.vim.lsp.omnifunc')
   end

   config.on_init = function(client, results)
      local buf_attach = function()
         vim.lsp.buf_attach_client(0, client.id)
      end

      autocmd = vim.api.nvim_create_autocmd('FileType', {
                                               desc = string.format('Attach LSP: %s', client.name),
                                               pattern = filetypes,
                                               callback = buf_attach
      })

      if vim.v.vim_did_enter == 1 and
         vim.tbl_contains(filetypes, vim.bo.filetype)
      then
         buf_attach()
      end
   end

   config.on_exit = vim.schedule_wrap(function(code, signal, client_id)
         vim.api.nvim_del_autocmd(autocmd)
   end)

   vim.lsp.start_client(config)
end


--
-- When an LSP attachment happens then ensure we bind
-- a key for hover-information
--
vim.api.nvim_create_autocmd('LspAttach', {
                               desc = 'LSP actions',
                               callback = function()
                                  local bufmap = function(mode, lhs, rhs)
                                     local opts = {buffer = true}
                                     vim.keymap.set(mode, lhs, rhs, opts)
                                  end

                                  -- Displays hover information about the symbol
                                  -- under the cursor when you press "K".
                                  bufmap('n', 'K', '<cmd>lua vim.lsp.buf.hover()<cr>')

                               end })


--
-- Invoke our helper function to start the server.
--
-- NOTE: This runs every time you launch neovim, even if you're not opening
--       a yal/lisp file.
--
launch_yal_server()
